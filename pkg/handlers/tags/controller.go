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
		tagsOnModeration.GET("", h.GetTagsPool)
		tagsOnModeration.POST("/:id", h.ApproveTag)
		tagsOnModeration.DELETE("/:id", h.DeclineTag)
	}
}
