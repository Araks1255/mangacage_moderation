package titles

import (
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB           *gorm.DB
	TitlesCovers *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, mongoClient *mongo.Client, secretKey string, r *gin.Engine) {
	titlesCoversCollection := mongoClient.Database("mangacage").Collection(mongodb.TitlesCoversCollection)

	h := handler{
		DB:           db,
		TitlesCovers: titlesCoversCollection,
	}

	titlesOnModeration := r.Group("/moderation/api/titles-on-moderation")
	{
		titlesOnModeration.GET("/", h.GetTitlesOnModerationPool)
		titlesOnModeration.POST("/:id", h.ApproveTitleOnModeration)
		titlesOnModeration.DELETE("/:id", h.DeclineTitleOnModeration)
		titlesOnModeration.PATCH("/:id/review", h.ReviewTitleOnModeration)
		titlesOnModeration.GET("/:id/cover", h.GetTitleOnModerationCover)
		titlesOnModeration.GET("/reviewing-by-me", h.GetTitlesOnModerationReviewingByMe)
	}

	titles := r.Group("/moderation/api/titles")
	{
		titles.PATCH("/:id/hide", h.HideTitle)
		titles.PATCH("/:id/unhide", h.UnhideTitle)
	}
}
