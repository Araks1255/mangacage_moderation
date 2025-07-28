package authors

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
