package users

import (
	"errors"
	"log"
	"strconv"

	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetUserProfileChangesProfilePicture(c *gin.Context) {
	userProfileChangesID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id изменений профиля"})
		return
	}

	var result mongoModels.UserOnModerationProfilePicture

	filter := bson.M{"user_on_moderation_id": userProfileChangesID}

	if err := h.UsersProfilePictures.FindOne(c.Request.Context(), filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.AbortWithStatusJSON(404, gin.H{"error": "не найдено аватарки для изменений профиля"})
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", result.ProfilePicture)
}
