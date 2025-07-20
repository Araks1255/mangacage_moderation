package users

import (
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                   *gorm.DB
	UsersProfilePictures *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, mongoClient *mongo.Client, secretKey string, r *gin.Engine) {
	usersProfilePicturesCollection := mongoClient.Database("mangacage").Collection(mongodb.UsersProfilePicturesCollection)

	h := handler{
		DB:                   db,
		UsersProfilePictures: usersProfilePicturesCollection,
	}

	users := r.Group("/moderation/api/users")
	{
		users.GET("", h.GetUsersPool)
		users.PATCH("/:id", h.VerificateUser)
		users.POST("/:id", h.ApproveUserProfileChanges)
		users.DELETE("/:id", h.DeclineUser)
	}
}
