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
		chaptersOnModeration.POST("/:id", h.ApproveChapterOnModeration)
		chaptersOnModeration.DELETE("/:id", h.DeclineChapterOnModeration)
		chaptersOnModeration.GET("/", h.GetChaptersOnModerationPool)
		chaptersOnModeration.PATCH("/:id/review", h.ReviewChapterOnModeration)
		chaptersOnModeration.GET("/:id/page/:page", h.GetChapterOnModerationPage)
		chaptersOnModeration.GET("/reviewing-by-me", h.GetChaptersOnModerationReviewingByMe)
		chaptersOnModeration.GET("/:id", h.GetChapterOnModeration)
	}

	chapters := r.Group("/moderation/api/chapters")
	{
		chapters.PATCH("/:id/hide", h.HideChapter)
		chapters.PATCH("/:id/unhide", h.UnhideChapter)
	}
}
