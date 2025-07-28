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
		usersOnVerification := users.Group("/on-verification")
		{
			usersOnVerification.GET("/", h.GetUsersOnVerificationPool)
			usersOnVerification.PATCH("/:id", h.VerificateUser)
			usersOnVerification.PATCH("/:id/review", h.ReviewUserVerification)
			usersOnVerification.GET("/reviewing-by-me", h.GetUsersOnVerificationReviewingByMe)
			usersOnVerification.GET("/:id", h.GetUserOnVerification)
		}

		profileChanges := users.Group("/profile-changes")
		{
			profileChanges.POST("/:id", h.ApproveUserProfileChanges)
			profileChanges.DELETE("/:id", h.DeclineUserProfileChanges)
			profileChanges.GET("/", h.GetUsersProfileChangesPool)
			profileChanges.PATCH("/:id/review", h.ReviewUserProfileChange)
			profileChanges.GET("/:id/profile-picture", h.GetUserProfileChangesProfilePicture)
			profileChanges.GET("/reviewing-by-me", h.GetUserProfileChangesReviewingByMe)
			profileChanges.GET("/:id", h.GetUserProfileChanges)
		}
	}
}
