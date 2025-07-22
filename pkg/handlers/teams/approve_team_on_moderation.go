package teams

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func (h handler) ApproveTeamOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	teamOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id команды на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	teamOnModeration, code, err := popTeamOnModeration(tx, uint(teamOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if teamOnModeration.ExistingID == nil {
		err := createTeam(c.Request.Context(), tx, h.TeamsCovers, *teamOnModeration)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	if teamOnModeration.ExistingID != nil {
		if err := updateTeam(c.Request.Context(), tx, h.TeamsCovers, *teamOnModeration); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "заявка на модерацию команды успешно одобрена"})
	// Уведомление
}

func popTeamOnModeration(db *gorm.DB, teamOnModerationID, userID uint) (team *models.TeamOnModeration, code int, err error) {
	var result models.TeamOnModeration

	err = db.Raw(
		"DELETE FROM teams_on_moderation WHERE id = ? AND moderator_id = ?",
		teamOnModerationID, userID,
	).Scan(&result).Error

	if err != nil {
		return nil, 500, err
	}

	if result.ID == 0 {
		return nil, 404, errors.New("команда на модерации не найдена среди заявок под вашим рассмотрением")
	}

	return &result, 0, nil
}

func createTeam(ctx context.Context, db *gorm.DB, collection *mongo.Collection, teamOnModeration models.TeamOnModeration) error {
	newTeamID, err := insertTeam(db, teamOnModeration.ToTeam())
	if err != nil {
		return err
	}

	if err := makeUserTeamLeader(db, teamOnModeration.CreatorID, newTeamID); err != nil {
		return err
	}

	if err := replaceTeamCoverTeamOnModerationID(ctx, collection, teamOnModeration.ID, newTeamID); err != nil {
		return err
	}

	return nil
}

func insertTeam(db *gorm.DB, team models.Team) (uint, error) {
	err := db.Create(&team).Error

	if err != nil {
		return 500, err
	}

	return team.ID, nil
}

func makeUserTeamLeader(db *gorm.DB, userID, teamID uint) error {
	err := db.Exec("UPDATE users SET team_id = ? WHERE id = ?", teamID, userID).Error

	if err != nil {
		return err
	}

	err = db.Exec("INSERT INTO user_roles (user_id, role_id) SELECT ? (SELECT id FROM roles WHERE name = 'team_leader')", userID).Error

	if err != nil {
		return err
	}

	return nil
}

func replaceTeamCoverTeamOnModerationID(ctx context.Context, collection *mongo.Collection, teamOnModerationID, teamID uint) error {
	filter := bson.M{"team_on_moderation_id": teamOnModerationID}

	update := bson.M{
		"$set":   bson.M{"team_id": teamID},
		"$unset": bson.M{"team_on_moderation_id": ""},
	}

	res, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("не удалось снять с модерации обложку команды")
	}

	return nil
}

func updateTeam(ctx context.Context, db *gorm.DB, collection *mongo.Collection, teamOnModeration models.TeamOnModeration) error {
	if err := db.Model("teams_on_moderation").Updates(&teamOnModeration).Error; err != nil {
		return err
	}

	if err := updateTeamCover(ctx, collection, teamOnModeration.ID, *teamOnModeration.ExistingID); err != nil {
		return err
	}

	return nil
}

func updateTeamCover(ctx context.Context, collection *mongo.Collection, teamOnModerationID, teamID uint) error {
	filter := bson.M{"team_id": teamID}
	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	filter = bson.M{"team_on_moderation_id": teamOnModerationID}
	update := bson.M{
		"$set":   bson.M{"team_id": teamID},
		"$unset": bson.M{"team_on_moderation_id": teamOnModerationID},
	}

	res, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return nil
	}

	return nil
}
