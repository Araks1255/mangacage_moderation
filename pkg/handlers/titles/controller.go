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
		titlesOnModeration.GET("/", h.GetTitlesPool)
		titlesOnModeration.POST("/:id", h.ApproveTitle)
		titlesOnModeration.DELETE("/:id", h.DeclineTitle)
	}
}
