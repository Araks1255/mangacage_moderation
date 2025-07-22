package chapters

import (
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB            *gorm.DB
	ChaptersPages *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, mongoClient *mongo.Client, secretKey string, r *gin.Engine) {
	chaptersPagesCollection := mongoClient.Database("mangacage").Collection(mongodb.ChaptersPagesCollection)

	h := handler{
		DB:            db,
		ChaptersPages: chaptersPagesCollection,
	}

	chaptersOnModeration := r.Group("/moderation/api/chapters-on-moderation")
	{
		chaptersOnModeration.POST("/:id", h.ApproveChapter)
		chaptersOnModeration.DELETE("/:id", h.DeclineChapter)
		chaptersOnModeration.GET("/", h.GetChaptersPool)
	}

	chapters := r.Group("/moderation/api/chapters")
	{
		chapters.PATCH("/:id/hide", h.HideChapter)
		chapters.PATCH("/:id/unhide", h.UnhideChapter)
	}
}
