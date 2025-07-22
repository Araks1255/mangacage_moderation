package teams

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func (h handler) DeclineTeamOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	teamOnModerationID, reason, err := parseDeclineTeamBody(c.ShouldBindJSON, c.Param)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	code, err := deleteTeamOnModeration(tx, teamOnModerationID, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := deleteTeamOnModerationCover(c.Request.Context(), h.TeamsCovers, teamOnModerationID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на модерацию команды успешно отклонена"})
	// Уведомление с причиной
	log.Println(reason)
}

func parseDeclineTeamBody(bindFn func(any) error, paramFn func(string) string) (teamID uint, reason string, err error) {
	var requestBody struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := bindFn(&requestBody); err != nil {
		return 0, "", err
	}

	id, err := strconv.ParseUint(paramFn("id"), 10, 64)
	if err != nil {
		return 0, "", err
	}

	return uint(id), requestBody.Reason, nil
}

func deleteTeamOnModeration(db *gorm.DB, teamOnModerationID, userID uint) (code int, err error) {
	result := db.Exec(
		"DELETE FROM teams_on_moderation WHERE id = ? AND moderator_id = ?",
		teamOnModerationID, userID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("команда на модерации не найдена среди рассматриваемых вами")
	}

	return 0, nil
}

func deleteTeamOnModerationCover(ctx context.Context, collection *mongo.Collection, teamOnModerationID uint) error {
	filter := bson.M{"team_on_moderation_id": teamOnModerationID}

	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
