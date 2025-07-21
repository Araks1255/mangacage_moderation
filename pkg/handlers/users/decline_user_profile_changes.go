package users

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

	code, err := deleteUserProfileChanges(tx, profileChangesID, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := deleteProfileChangesCover(c.Request.Context(), h.UsersProfilePictures, profileChangesID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на модерацию изменений профиля успешно отклонена"})
	// Уведомление с причиной
	log.Println(reason)
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

func deleteUserProfileChanges(db *gorm.DB, profileChangesID, userID uint) (code int, err error) {
	result := db.Exec(
		"DELETE FROM users_on_moderation WHERE id = ? AND moderator_id = ?",
		profileChangesID, userID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("изменения профиля пользователя не найдены среди рассматриваемых вами")
	}

	return 0, nil
}

func deleteProfileChangesCover(ctx context.Context, collection *mongo.Collection, profileChangesID uint) error {
	filter := bson.M{"user_on_moderation_id": profileChangesID}

	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
