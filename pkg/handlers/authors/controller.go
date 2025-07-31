package authors

import (
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	NotificationsClient pb.ModerationNotificationsClient
}

func RegisterRoutes(db *gorm.DB, secretKey string, notificationsClient pb.ModerationNotificationsClient, r *gin.Engine) {
	h := handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}

	authorsOnModeration := r.Group("/moderation/api/authors-on-moderation")
	{
		authorsOnModeration.GET("", h.GetAuthorsOnModerationPool)
		authorsOnModeration.POST("/:id", h.ApproveAuthorOnModeration)
		authorsOnModeration.DELETE("/:id", h.DeclineAuthorOnModeration)
		authorsOnModeration.PATCH("/:id/review", h.ReviewAuthorOnModeration)
		authorsOnModeration.GET("/reviewing-by-me", h.GetAuthorsOnModerationReviewingByMe)
		authorsOnModeration.GET("/:id", h.GetAuthorOnModeration)
	}
}
