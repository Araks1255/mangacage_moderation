package teams

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
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

	teamOnModeration, code, err := deleteTeamOnModeration(tx, teamOnModerationID, claims.ID)
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

	var name string
	if teamOnModeration.Name != nil {
		name = *teamOnModeration.Name
	}

	if _, err := h.NotificationsClient.SendModerationRequestDeclineReason(
		c.Request.Context(), &pb.ModerationRequestDeclineReason{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_CHAPTER,
			EntityName:         name,
			CreatorID:          uint64(teamOnModeration.CreatorID),
			Reason:             reason,
		},
	); err != nil {
		log.Println(err)
	}
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

func deleteTeamOnModeration(db *gorm.DB, teamOnModerationID, userID uint) (deleted *models.TeamOnModeration, code int, err error) {
	var deletedTeamOnModeration models.TeamOnModeration

	err = db.Raw(
		"DELETE FROM teams_on_moderation WHERE id = ? AND moderator_id = ? RETURNING id, name, creator_id",
		teamOnModerationID, userID,
	).Scan(&deletedTeamOnModeration).Error

	if err != nil {
		return nil, 500, err
	}

	if deletedTeamOnModeration.ID == 0 {
		return nil, 404, errors.New("команда на модерации не найдена среди рассматриваемых вами")
	}

	return &deletedTeamOnModeration, 0, nil
}

func deleteTeamOnModerationCover(ctx context.Context, collection *mongo.Collection, teamOnModerationID uint) error {
	filter := bson.M{"team_on_moderation_id": teamOnModerationID}

	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
