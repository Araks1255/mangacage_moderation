package teams

import (
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB          *gorm.DB
	TeamsCovers *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, mongoClient *mongo.Client, secretKey string, r *gin.Engine) {
	teamsCoversCollection := mongoClient.Database("mangacage").Collection(mongodb.TeamsCoversCollection)

	h := handler{
		DB:          db,
		TeamsCovers: teamsCoversCollection,
	}

	teamsOnModeration := r.Group("/moderation/api/teams-on-moderation")
	{
		teamsOnModeration.GET("", h.GetTeamsPool)
		teamsOnModeration.POST("/:id", h.ApproveTeam)
		teamsOnModeration.DELETE("/:id", h.DeclineTeam)
	}
}
