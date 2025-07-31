package tags

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

	tagsOnModeration := r.Group("/moderation/api/tags-on-moderation")
	{
		tagsOnModeration.GET("", h.GetTagsOnModerationPool)
		tagsOnModeration.POST("/:id", h.ApproveTagOnModeration)
		tagsOnModeration.DELETE("/:id", h.DeclineTagOnModeration)
		tagsOnModeration.PATCH("/:id/review", h.ReviewTagOnModeration)
		tagsOnModeration.GET("/reviewing-by-me", h.GetTagsOnModerationReviewingByMe)
		tagsOnModeration.GET("/:id", h.GetTagOnModeration)
	}
}
