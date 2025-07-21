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
		genresOnModeration.GET("", h.GetGenresPool)
		genresOnModeration.POST("/:id", h.ApproveGenre)
		genresOnModeration.DELETE("/:id", h.DeclineGenre)
	}
}
