package genres

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, secretKey string, r *gin.Engine) {
	h := handler{
		DB: db,
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
