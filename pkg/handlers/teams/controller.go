package teams

import (
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	TeamsCovers         *mongo.Collection
	NotificationsClient pb.ModerationNotificationsClient
}

func RegisterRoutes(db *gorm.DB, mongoClient *mongo.Client, notificationsClient pb.ModerationNotificationsClient, secretKey string, r *gin.Engine) {
	teamsCoversCollection := mongoClient.Database("mangacage").Collection(mongodb.TeamsCoversCollection)

	h := handler{
		DB:                  db,
		TeamsCovers:         teamsCoversCollection,
		NotificationsClient: notificationsClient,
	}

	teamsOnModeration := r.Group("/moderation/api/teams-on-moderation")
	{
		teamsOnModeration.GET("", h.GetTeamsOnModerationPool)
		teamsOnModeration.POST("/:id", h.ApproveTeamOnModeration)
		teamsOnModeration.DELETE("/:id", h.DeclineTeamOnModeration)
		teamsOnModeration.PATCH("/:id/review", h.ReviewTeamOnModeration)
		teamsOnModeration.GET("/:id/cover", h.GetTeamOnModerationCover)
		teamsOnModeration.GET("/reviewing-by-me", h.GetTeamsOnModerationReviewingByMe)
		teamsOnModeration.GET("/:id", h.GetTeamOnModeration)
	}
}
