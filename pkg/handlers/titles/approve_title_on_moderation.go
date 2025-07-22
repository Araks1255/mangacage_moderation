package titles

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/helpers/titles"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	if titleOnModeration.ExistingID == nil {
		err = createTitle(c.Request.Context(), tx, h.TitlesCovers, *titleOnModeration)
	} else {
		err = updateTitle(c.Request.Context(), tx, h.TitlesCovers, *titleOnModeration)
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if _, err := titles.DeleteTitleOnModeration(tx, uint(titleOnModerationID), claims.ID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на модерацию тайтла успешно одобрена"})
	// Уведомление
}

func getTitleOnModeration(db *gorm.DB, titleOnModerationID, userID uint) (chapter *models.TitleOnModeration, code int, err error) {
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

func createTitle(ctx context.Context, db *gorm.DB, collection *mongo.Collection, titleOnModeration models.TitleOnModeration) error {
	newTitleID, err := insertTitle(db, titleOnModeration.ToTitle())

	if err != nil {
		return err
	}

	err = replaceTitleCoverTitleOnModerationID(ctx, collection, titleOnModeration.ID, newTitleID)

	if err != nil {
		return err
	}

	err = replaceChaptersOnModerationTitleOnModerationID(db, titleOnModeration.ID, newTitleID)
	if err != nil {
		return err
	}

	return nil
}

func updateTitle(ctx context.Context, db *gorm.DB, collection *mongo.Collection, titleOnModeration models.TitleOnModeration) error {
	title := titleOnModeration.ToTitle()

	if err := db.Model("titles_on_moderation").Updates(&title).Error; err != nil {
		return err
	}

	if err := updateTitleCover(ctx, collection, titleOnModeration.ID, *titleOnModeration.ExistingID); err != nil {
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

func updateTitleCover(ctx context.Context, collection *mongo.Collection, titleOnModerationID, titleID uint) error {
	filter := bson.M{"title_id": titleID}
	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	filter = bson.M{"title_on_moderation_id": titleOnModerationID}
	update := bson.M{
		"$set":   bson.M{"title_id": titleID},
		"$unset": bson.M{"title_on_moderation_id": titleOnModerationID},
	}

	res, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return nil
	}

	return nil
}
