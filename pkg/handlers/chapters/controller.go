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

	chapters := r.Group("/moderation/api/chapters")
	{
		chapters.POST("/:id", h.ApproveChapter)
		chapters.DELETE("/:id", h.DeclineChapter)
		chapters.GET("/", h.GetChaptersPool)
	}
}
