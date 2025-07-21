package main

import (
	"context"

	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage_moderation/pkg/common/db"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/authors"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/genres"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/tags"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/teams"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/titles"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/users"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	ctx := context.Background()

	dbUrl := viper.Get("DB_URL").(string)
	mongoUrl := viper.Get("MONGO_URL").(string)
	secretKey := viper.Get("SECRET_KEY").(string)

	mongoClient, err := db.MongoInit(mongoUrl)
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(ctx)

	db, err := db.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(middlewares.Auth(secretKey), middlewares.RequireRoles(db, []string{"moderator", "admin"}))

	chapters.RegisterRoutes(db, mongoClient, secretKey, router)
	titles.RegisterRoutes(db, mongoClient, secretKey, router)
	teams.RegisterRoutes(db, mongoClient, secretKey, router)
	users.RegisterRoutes(db, mongoClient, secretKey, router)
	tags.RegisterRoutes(db, secretKey, router)
	genres.RegisterRoutes(db, secretKey, router)
	authors.RegisterRoutes(db, secretKey, router)

	router.Run(":8080")
}
