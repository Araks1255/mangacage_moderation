package genres

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

	genresOnModeration := r.Group("/moderation/api/genres-on-moderation")
	{
		genresOnModeration.GET("", h.GetGenresOnModerationPool)
		genresOnModeration.POST("/:id", h.ApproveGenreOnModeration)
		genresOnModeration.DELETE("/:id", h.DeclineGenreOnModeration)
		genresOnModeration.PATCH("/:id/review", h.ReviewGenreOnModeration)
		genresOnModeration.GET("/reviewing-by-me", h.GetGenresOnModerationReviewingByMe)
		genresOnModeration.GET("/:id", h.GetGenreOnModeration)
	}
}
