package users

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

func (h handler) ApproveUserProfileChanges(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	userOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id изменений профиля пользователя"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	userOnModeration, code, err := popUserOnModeration(tx, uint(userOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := updateUser(c.Request.Context(), tx, h.UsersProfilePictures, *userOnModeration); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения профиля пользователя успешно одобрены"})

	if _, err := h.NotificationsClient.NotifyAboutApprovedModerationRequest(
		c.Request.Context(), &pb.ApprovedEntity{
			Entity:    enums.Entity_ENTITY_PROFILE,
			CreatorID: uint64(*userOnModeration.ExistingID),
		},
	); err != nil {
		log.Println(err)
	}
}

func popUserOnModeration(db *gorm.DB, userOnModerationID, moderatorID uint) (userOnModeration *models.UserOnModeration, code int, err error) {
	var result models.UserOnModeration

	err = db.Raw(
		"DELETE FROM users_on_moderation WHERE id = ? AND moderator_id = ? RETURNING *",
		userOnModerationID, moderatorID,
	).Scan(&result).Error

	if err != nil {
		return nil, 500, err
	}

	if result.ID == 0 {
		return nil, 404, errors.New("изменения профиля пользователя не найдены среди рассматриваемых вами")
	}

	return &result, 0, nil
}

func updateUser(ctx context.Context, db *gorm.DB, collection *mongo.Collection, userOnModeration models.UserOnModeration) error {
	user := userOnModeration.ToUser()
	if err := db.Table("users").Where("id = ?", userOnModeration.ExistingID).Updates(&user).Error; err != nil {
		return err
	}

	if err := updateUserProfilePicture(ctx, collection, userOnModeration.ID, *userOnModeration.ExistingID); err != nil {
		return err
	}

	return nil
}

func updateUserProfilePicture(ctx context.Context, collection *mongo.Collection, userOnModerationID, userID uint) error {
	userOnModerationFilter := bson.M{"user_on_moderation_id": userOnModerationID}

	var profilePicture mongoModels.UserOnModerationProfilePicture

	if err := collection.FindOne(ctx, userOnModerationFilter).Decode(&profilePicture); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}

	userFilter := bson.M{"user_id": userID}
	userUpdate := bson.M{"$set": bson.M{"profile_picture": profilePicture.ProfilePicture}}
	userOpts := options.Update().SetUpsert(true) // на всякий случай

	res, err := collection.UpdateOne(ctx, userFilter, userUpdate, userOpts)

	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		log.Printf("не удалось изменить аватарку пользователя.\nid пользователя: %d\nid изменений профиля: %d", userID, userOnModerationID)
	}

	delteRes, err := collection.DeleteOne(ctx, userOnModerationFilter)

	if err != nil {
		log.Printf("не удалось удалить аватарку изменений профиля.\nошибка: %s\nid: %d", err.Error(), userOnModerationID)
	}

	if delteRes.DeletedCount == 0 {
		log.Printf("не удалось удалить аватарку изменений профиля. id: %d", userOnModerationID)
	}

	return nil
}
