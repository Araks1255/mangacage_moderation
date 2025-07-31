package titles

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/helpers/titles"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

func (h handler) ApproveTitleOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	titleOnModeration, code, err := getTitleOnModeration(tx, uint(titleOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	var titleID uint

	if titleOnModeration.ExistingID == nil {
		titleID, err = createTitle(c.Request.Context(), tx, h.TitlesCovers, *titleOnModeration)
	} else {
		err = updateTitle(c.Request.Context(), tx, h.TitlesCovers, *titleOnModeration)
		titleID = *titleOnModeration.ExistingID
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if _, _, err := titles.DeleteTitleOnModeration(tx, uint(titleOnModerationID), claims.ID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на модерацию тайтла успешно одобрена"})

	if _, err := h.NotificationsClient.NotifyAboutApprovedModerationRequest(
		c.Request.Context(), &pb.ApprovedEntity{
			Entity:    enums.Entity_ENTITY_TITLE,
			ID:        uint64(titleID),
			CreatorID: uint64(titleOnModeration.CreatorID),
		},
	); err != nil {
		log.Println(err)
	}
}

func getTitleOnModeration(db *gorm.DB, titleOnModerationID, userID uint) (title *models.TitleOnModeration, code int, err error) {
	var result models.TitleOnModeration

	err = db.Raw(
		"SELECT * FROM titles_on_moderation WHERE id = ? AND moderator_id = ?",
		titleOnModerationID, userID,
	).Scan(&result).Error

	if err != nil {
		return nil, 500, err
	}

	if result.ID == 0 {
		return nil, 404, errors.New("тайтл на модерации не найден среди заявок под вашим рассмотрением")
	}

	if result.AuthorOnModerationID != nil {
		return nil, 409, errors.New("для начала необходимо снять автора тайтла с модерации")
	}

	return &result, 0, nil
}

func createTitle(ctx context.Context, db *gorm.DB, collection *mongo.Collection, titleOnModeration models.TitleOnModeration) (uint, error) {
	newTitleID, err := insertTitle(db, titleOnModeration.ToTitle())
	if err != nil {
		return 0, err
	}

	if err = setTitleGenresFromTitleOnModeration(db, newTitleID, titleOnModeration.ID); err != nil {
		return 0, err
	}

	if err := setTitleTagsFromTitleOnModeration(db, newTitleID, titleOnModeration.ID); err != nil {
		return 0, err
	}

	if err = replaceChaptersOnModerationTitleOnModerationID(db, titleOnModeration.ID, newTitleID); err != nil {
		return 0, err
	}

	if err := makeTitleTranslatingByCreatorTeam(db, newTitleID, titleOnModeration.CreatorID); err != nil {
		return 0, err
	}

	if err = replaceTitleCoverTitleOnModerationID(ctx, collection, titleOnModeration.ID, newTitleID); err != nil {
		return 0, err
	}

	return newTitleID, nil
}

func updateTitle(ctx context.Context, db *gorm.DB, collection *mongo.Collection, titleOnModeration models.TitleOnModeration) error {
	title := titleOnModeration.ToTitle()

	if err := db.Table("titles").Where("id = ?", titleOnModeration.ExistingID).Updates(&title).Error; err != nil {
		return err
	}

	if err := updateTitleCover(ctx, collection, titleOnModeration.ID, *titleOnModeration.ExistingID, titleOnModeration.CreatorID); err != nil {
		return err
	}

	return nil
}

func insertTitle(db *gorm.DB, title models.Title) (uint, error) {
	err := db.Create(&title).Error

	if err != nil {
		return 500, err
	}

	return title.ID, nil
}

func replaceTitleCoverTitleOnModerationID(ctx context.Context, collection *mongo.Collection, titleOnModerationID, titleID uint) error {
	filter := bson.M{"title_on_moderation_id": titleOnModerationID}

	update := bson.M{
		"$set":   bson.M{"title_id": titleID},
		"$unset": bson.M{"title_on_moderation_id": ""},
	}

	res, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("не удалось снять с модерации обложку тайтла")
	}

	return nil
}

func replaceChaptersOnModerationTitleOnModerationID(db *gorm.DB, titleOnModerationID, titleID uint) error {
	return db.Exec(
		`UPDATE chapters_on_moderation SET
			title_id = ?,
			title_on_moderation_id = NULL
		WHERE
			title_on_moderation_id = ?`,
		titleID, titleOnModerationID,
	).Error
}

func makeTitleTranslatingByCreatorTeam(db *gorm.DB, titleID, creatorID uint) error {
	return db.Exec(
		`INSERT INTO title_teams
			(title_id, team_id)
		SELECT
			?, teams.id
		FROM
			users AS u
			INNER JOIN teams ON teams.id = u.team_id
			INNER JOIN user_roles AS ur ON ur.user_id = u.id
			INNER JOIN roles AS r ON r.id = ur.role_id
		WHERE
			u.id = ? AND r.name = 'team_leader'`,
		titleID, creatorID,
	).Error
}

func updateTitleCover(ctx context.Context, collection *mongo.Collection, titleOnModerationID, titleID, creatorID uint) error {
	titleOnModerationFilter := bson.M{"title_on_moderation_id": titleOnModerationID}

	var cover mongoModels.TitleOnModerationCover

	if err := collection.FindOne(ctx, titleOnModerationFilter).Decode(&cover); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}

	titleFilter := bson.M{"title_id": titleID}
	titleUpdate := bson.M{"$set": bson.M{"cover": cover.Cover, "creator_id": creatorID}}
	titleOpts := options.Update().SetUpsert(true) // на всякий случай

	res, err := collection.UpdateOne(ctx, titleFilter, titleUpdate, titleOpts)

	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		log.Printf("не удалось изменить обложку тайтла.\nid тайтла: %d\nid тайтла на модерации: %d", titleID, titleOnModerationID)
	}

	delteRes, err := collection.DeleteOne(ctx, titleOnModerationFilter)

	if err != nil {
		log.Printf("не удалось удалить обложку тайтла на модерации.\nошибка: %s\nid: %d", err.Error(), titleOnModerationID)
	}

	if delteRes.DeletedCount == 0 {
		log.Printf("не удалось удалить обложку тайтла на модерации. id: %d", titleOnModerationID)
	}

	return nil
}

func setTitleGenresFromTitleOnModeration(db *gorm.DB, titleID, titleOnModerationID uint) error {
	var genresExist bool

	err := db.Raw(
		"SELECT EXISTS(SELECT 1 FROM title_on_moderation_genres WHERE title_on_moderation_id = ?)",
		titleOnModerationID,
	).Scan(&genresExist).Error

	if err != nil {
		return err
	}

	if !genresExist {
		return nil
	}

	if err := db.Exec("DELETE FROM title_genres WHERE title_id = ?", titleID).Error; err != nil {
		return err
	}

	result := db.Exec(
		`INSERT INTO
			title_genres (title_id, genre_id)
		SELECT
			?, genre_id
		FROM
			title_on_moderation_genres
		WHERE
			title_on_moderation_id = ?`,
		titleID, titleOnModerationID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(
			"не удалось перенести жанры тайтла на модерации созданному тайтлу.\nid тайтла на модерации - %d",
			titleOnModerationID,
		)
	}

	return nil
}

func setTitleTagsFromTitleOnModeration(db *gorm.DB, titleID, titleOnModerationID uint) error {
	var tagsExist bool

	err := db.Raw(
		"SELECT EXISTS(SELECT 1 FROM title_on_moderation_tags WHERE title_on_moderation_id = ?)",
		titleOnModerationID,
	).Scan(&tagsExist).Error

	if err != nil {
		return err
	}

	if !tagsExist {
		return nil
	}

	if err := db.Exec("DELETE FROM title_tags WHERE title_id = ?", titleID).Error; err != nil {
		return err
	}

	result := db.Exec(
		`INSERT INTO
			title_tags (title_id, tag_id)
		SELECT
			?, tag_id
		FROM
			title_on_moderation_tags
		WHERE
			title_on_moderation_id = ?`,
		titleID, titleOnModerationID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(
			"не удалось перенести теги тайтла на модерации созданному тайтлу.\nid тайтла на модерации - %d",
			titleOnModerationID,
		)

	}

	return nil
}
