package teams

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	var teamID uint

	if teamOnModeration.ExistingID == nil {
		teamID, err = createTeam(c.Request.Context(), tx, h.TeamsCovers, *teamOnModeration)
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
		teamID = *teamOnModeration.ExistingID
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "заявка на модерацию команды успешно одобрена"})

	if _, err := h.NotificationsClient.NotifyAboutApprovedModerationRequest(
		c.Request.Context(), &pb.ApprovedEntity{
			Entity:    enums.Entity_ENTITY_TEAM,
			ID:        uint64(teamID),
			CreatorID: uint64(*teamOnModeration.ExistingID),
		},
	); err != nil {
		log.Println(err)
	}
}

func popTeamOnModeration(db *gorm.DB, teamOnModerationID, userID uint) (team *models.TeamOnModeration, code int, err error) {
	var result models.TeamOnModeration

	err = db.Raw(
		"DELETE FROM teams_on_moderation WHERE id = ? AND moderator_id = ? RETURNING *",
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

func createTeam(ctx context.Context, db *gorm.DB, collection *mongo.Collection, teamOnModeration models.TeamOnModeration) (uint, error) {
	newTeamID, err := insertTeam(db, teamOnModeration.ToTeam())
	if err != nil {
		return 0, err
	}

	if err := makeUserTeamLeader(db, teamOnModeration.CreatorID, newTeamID); err != nil {
		return 0, err
	}

	if err := replaceTeamCoverTeamOnModerationID(ctx, collection, teamOnModeration.ID, newTeamID); err != nil {
		return 0, err
	}

	return newTeamID, nil
}

func updateTeam(ctx context.Context, db *gorm.DB, collection *mongo.Collection, teamOnModeration models.TeamOnModeration) error {
	team := teamOnModeration.ToTeam()

	if err := db.Table("teams").Where("id = ?", teamOnModeration.ExistingID).Updates(&team).Error; err != nil {
		return err
	}

	if err := updateTeamCover(ctx, collection, teamOnModeration.ID, *teamOnModeration.ExistingID, teamOnModeration.CreatorID); err != nil {
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

	err = db.Exec("INSERT INTO user_roles (user_id, role_id) SELECT ?, (SELECT id FROM roles WHERE name = 'team_leader')", userID).Error

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

func updateTeamCover(ctx context.Context, collection *mongo.Collection, teamOnModerationID, teamID, creatorID uint) error {
	teamOnModerationFilter := bson.M{"team_on_moderation_id": teamOnModerationID}

	var cover mongoModels.TeamOnModerationCover

	if err := collection.FindOne(ctx, teamOnModerationFilter).Decode(cover); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}

	teamFilter := bson.M{"team_id": teamID}
	teamUpdate := bson.M{"$set": bson.M{"cover": cover.Cover, "creator_id": creatorID}}
	teamOpts := options.Update().SetUpsert(true)

	res, err := collection.UpdateOne(ctx, teamFilter, teamUpdate, teamOpts)

	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		log.Printf("не удалось обновить обложку команды.\nid команды: %d\nid команды на модерации: %d", teamID, teamOnModerationID)
	}

	deleteRes, err := collection.DeleteOne(ctx, teamOnModerationFilter)

	if err != nil {
		log.Printf("не удалось удалить обложку команды на модерации.\nошибка: %s\nid: %d", err.Error(), teamOnModerationID)
	}

	if deleteRes.DeletedCount == 0 {
		log.Printf("не удалось удалить обложку команды на модерации.\nid: %d", teamOnModerationID)
	}

	return nil
}
