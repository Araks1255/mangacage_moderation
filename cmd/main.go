package main

import (
	"context"
	"flag"

	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage_moderation/pkg/common/db"
	"github.com/Araks1255/mangacage_moderation/pkg/common/seeder"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/authors"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/genres"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/tags"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/teams"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/titles"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/users"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/views"
	moderationMiddlewares "github.com/Araks1255/mangacage_moderation/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	ctx := context.Background()

	dbUrl := viper.Get("DB_URL").(string)
	mongoUrl := viper.Get("MONGO_URL").(string)
	secretKey := viper.Get("SECRET_KEY").(string)

	mongoClient, err := db.MongoInit(ctx, mongoUrl)
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(ctx)

	db, err := db.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	seedFlag := flag.Bool("seed", false, "")

	flag.Parse()

	if *seedFlag {
		if err := seeder.SeedEntitiesOnModeration(db, mongoClient.Database("mangacage")); err != nil {
			panic(err)
		}
	}

	modersIDs, err := getModersIDs(db)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(middlewares.Auth(secretKey), moderationMiddlewares.RequireModer(modersIDs))

	chapters.RegisterRoutes(db, mongoClient, secretKey, router)
	titles.RegisterRoutes(db, mongoClient, secretKey, router)
	teams.RegisterRoutes(db, mongoClient, secretKey, router)
	users.RegisterRoutes(db, mongoClient, secretKey, router)
	tags.RegisterRoutes(db, secretKey, router)
	genres.RegisterRoutes(db, secretKey, router)
	authors.RegisterRoutes(db, secretKey, router)
	views.RegisterRoutes(router)

	router.Run(":80")
}

func getModersIDs(db *gorm.DB) (map[uint]struct{}, error) {
	var ids []uint

	err := db.Raw(
		`SELECT DISTINCT
			u.id
		FROM
			users AS u
			INNER JOIN user_roles AS ur ON ur.user_id = u.id
			INNER JOIN roles AS r ON r.id = ur.role_id
		WHERE
			r.name = 'moder' OR r.name = 'admin'`,
	).Scan(&ids).Error

	if err != nil {
		return nil, err
	}

	res := make(map[uint]struct{})

	for i := 0; i < len(ids); i++ {
		res[ids[i]] = struct{}{}
	}

	return res, nil
}
