package titles

import (
	"context"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/helpers/titles"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) DeclineTitleOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleOnModerationID, reason, err := parseDeclineTitleBody(c.ShouldBindJSON, c.Param)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	titleOnModeration, code, err := titles.DeleteTitleOnModeration(tx, titleOnModerationID, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := deleteTitleCover(c.Request.Context(), h.TitlesCovers, titleOnModerationID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на модерацию тайтла успешно отклонена"})

	if _, err := h.NotificationsClient.SendModerationRequestDeclineReason(
		c.Request.Context(), &pb.ModerationRequestDeclineReason{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_TITLE,
			EntityName:         *titleOnModeration.Name,
			CreatorID:          uint64(titleOnModeration.CreatorID),
			Reason:             reason,
		},
	); err != nil {
		log.Println(err)
	}
}

func parseDeclineTitleBody(bindFn func(any) error, paramFn func(string) string) (titleID uint, reason string, err error) {
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

func deleteTitleCover(ctx context.Context, collection *mongo.Collection, titleOnModerationID uint) error {
	filter := bson.M{"title_on_moderation_id": titleOnModerationID}

	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
