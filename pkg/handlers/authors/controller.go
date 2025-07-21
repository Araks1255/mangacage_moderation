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
		authorsOnModeration.GET("", h.GetAuthorsPool)
		authorsOnModeration.POST("/:id", h.ApproveAuthor)
		authorsOnModeration.DELETE("/:id", h.DeclineAuthor)
	}
}
