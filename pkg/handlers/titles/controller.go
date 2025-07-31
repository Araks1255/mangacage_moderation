package titles

import (
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	TitlesCovers        *mongo.Collection
	NotificationsClient pb.ModerationNotificationsClient
}

func RegisterRoutes(db *gorm.DB, mongoClient *mongo.Client, notificationsClient pb.ModerationNotificationsClient, secretKey string, r *gin.Engine) {
	titlesCoversCollection := mongoClient.Database("mangacage").Collection(mongodb.TitlesCoversCollection)

	h := handler{
		DB:                  db,
		TitlesCovers:        titlesCoversCollection,
		NotificationsClient: notificationsClient,
	}

	titlesOnModeration := r.Group("/moderation/api/titles-on-moderation")
	{
		titlesOnModeration.GET("/", h.GetTitlesOnModerationPool)
		titlesOnModeration.POST("/:id", h.ApproveTitleOnModeration)
		titlesOnModeration.DELETE("/:id", h.DeclineTitleOnModeration)
		titlesOnModeration.PATCH("/:id/review", h.ReviewTitleOnModeration)
		titlesOnModeration.GET("/:id/cover", h.GetTitleOnModerationCover)
		titlesOnModeration.GET("/reviewing-by-me", h.GetTitlesOnModerationReviewingByMe)
		titlesOnModeration.GET("/:id", h.GetTitleOnModeraiton)
	}

	titles := r.Group("/moderation/api/titles")
	{
		titles.PATCH("/:id/hide", h.HideTitle)
		titles.PATCH("/:id/unhide", h.UnhideTitle)
	}
}
