package users

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func (h handler) DeclineUserProfileChanges(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	profileChangesID, reason, err := parseDeclineUserProfileChangesBody(c.ShouldBindJSON, c.Param)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	userID, code, err := deleteUserProfileChanges(tx, profileChangesID, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := deleteProfileChangesProfilePicture(c.Request.Context(), h.UsersProfilePictures, profileChangesID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на модерацию изменений профиля успешно отклонена"})

	if _, err := h.NotificationsClient.SendModerationRequestDeclineReason(
		c.Request.Context(), &pb.ModerationRequestDeclineReason{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_PROFILE_CHANGES,
			CreatorID:          uint64(userID),
			Reason:             reason,
		},
	); err != nil {
		log.Println(err)
	}
}

func parseDeclineUserProfileChangesBody(bindFn func(any) error, paramFn func(string) string) (profileChangesID uint, reason string, err error) {
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

func deleteUserProfileChanges(db *gorm.DB, profileChangesID, moderatorID uint) (userID uint, code int, err error) {
	var userIDPtr *uint

	err = db.Raw(
		"DELETE FROM users_on_moderation WHERE id = ? AND moderator_id = ? RETURNING existing_id",
		profileChangesID, moderatorID,
	).Scan(&userIDPtr).Error

	if err != nil {
		return 0, 500, err
	}

	if userIDPtr == nil {
		return 0, 404, errors.New("изменения профиля пользователя не найдены среди рассматриваемых вами")
	}

	return *userIDPtr, 0, nil
}

func deleteProfileChangesProfilePicture(ctx context.Context, collection *mongo.Collection, profileChangesID uint) error {
	filter := bson.M{"user_on_moderation_id": profileChangesID}

	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
