package tags

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) DeclineTagOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	tagOnModerationID, reason, err := parseDeclineTagBody(c.ShouldBindJSON, c.Param)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tagOnModeration, code, err := deleteTagOnModeration(h.DB, tagOnModerationID, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на модерацию тега успешно отклонена"})

	if _, err := h.NotificationsClient.SendModerationRequestDeclineReason(
		c.Request.Context(), &pb.ModerationRequestDeclineReason{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_TAG,
			EntityName:         tagOnModeration.Name,
			CreatorID:          uint64(tagOnModeration.CreatorID),
			Reason:             reason,
		},
	); err != nil {
		log.Println(err)
	}
}

func parseDeclineTagBody(bindFn func(any) error, paramFn func(string) string) (tagID uint, reason string, err error) {
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

func deleteTagOnModeration(db *gorm.DB, tagOnModerationID, userID uint) (deleted *models.TagOnModeration, code int, err error) {
	var deletedTagOnModeration models.TagOnModeration

	err = db.Raw(
		"DELETE FROM tags_on_moderation WHERE id = ? AND moderator_id = ? RETURNING name, creator_id",
		tagOnModerationID, userID,
	).Scan(&deletedTagOnModeration).Error

	if err != nil {
		return nil, 500, err
	}

	if deletedTagOnModeration.Name == "" {
		return nil, 404, errors.New("тег на модерации не найден среди рассматриваемых вами")
	}

	return &deletedTagOnModeration, 0, nil
}
