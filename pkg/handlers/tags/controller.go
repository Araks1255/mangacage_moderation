package tags

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

	tagsOnModeration := r.Group("/moderation/api/tags-on-moderation")
	{
		tagsOnModeration.GET("", h.GetTagsOnModerationPool)
		tagsOnModeration.POST("/:id", h.ApproveTagOnModeration)
		tagsOnModeration.DELETE("/:id", h.DeclineTagOnModeration)
		tagsOnModeration.PATCH("/:id/review", h.ReviewTagOnModeration)
		tagsOnModeration.GET("/reviewing-by-me", h.GetTagsOnModerationReviewingByMe)
	}
}
