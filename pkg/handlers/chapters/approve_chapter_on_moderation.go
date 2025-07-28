package chapters

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func (h handler) ApproveChapterOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapterOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	chapterOnModeration, code, err := popChapterOnModeration(tx, uint(chapterOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if chapterOnModeration.ExistingID == nil {
		newChapterID, err := insertChapter(tx, chapterOnModeration.ToChapter())
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		err = replaceChapterPagesChapterOnModerationID(c.Request.Context(), h.ChaptersPages, chapterOnModeration.ID, newChapterID, chapterOnModeration.CreatorID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	if chapterOnModeration.ExistingID != nil {
		if err := updateChapter(tx, chapterOnModeration.ToChapter()); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на модерацию главы успешно одобрена"})
	// Уведомление
}

func popChapterOnModeration(db *gorm.DB, chapterOnModerationID, userID uint) (chapter *models.ChapterOnModeration, code int, err error) {
	var result models.ChapterOnModeration

	err = db.Raw(
		"DELETE FROM chapters_on_moderation WHERE id = ? AND moderator_id = ? RETURNING *",
		chapterOnModerationID, userID,
	).Scan(&result).Error

	if err != nil {
		return nil, 500, err
	}

	if result.ID == 0 {
		return nil, 404, errors.New("глава на модерации не найдена среди заявок под вашим рассмотрением")
	}

	if result.TitleOnModerationID != nil {
		return nil, 409, errors.New("для начала необходимо снять тайтл с модерации")
	}

	return &result, 0, nil
}

func insertChapter(db *gorm.DB, chapter models.Chapter) (uint, error) {
	err := db.Create(&chapter).Error

	if err != nil {
		return 500, err
	}

	return chapter.ID, nil
}

func updateChapter(db *gorm.DB, chapter models.Chapter) error {
	return db.Model("chapters_on_moderation").Updates(&chapter).Error
}

func replaceChapterPagesChapterOnModerationID(ctx context.Context, collection *mongo.Collection, chapterOnModerationID, chapterID, creatorID uint) error {
	filter := bson.M{"chapter_on_moderation_id": chapterOnModerationID}

	update := bson.M{
		"$set":   bson.M{"chapter_id": chapterID, "creator_id": creatorID},
		"$unset": bson.M{"chapter_on_moderation_id": ""},
	}

	res, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("не удалось снять с модерации страницы главы")
	}

	return nil
}
