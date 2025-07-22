package tags

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) ApproveTagOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	tagOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тега на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	tagOnModeration, code, err := popTagOnModeration(tx, uint(tagOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := createTag(tx, tagOnModeration.ToTag()); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "заявка на модерацию тега успешно одобрена"})
	// Уведомление
}

func popTagOnModeration(db *gorm.DB, tagOnModerationID, userID uint) (tag *models.TagOnModeration, code int, err error) {
	var result models.TagOnModeration

	err = db.Raw(
		"DELETE FROM tags_on_moderation WHERE id = ? AND moderator_id = ?",
		tagOnModerationID, userID,
	).Scan(&result).Error

	if err != nil {
		return nil, 500, err
	}

	if result.ID == 0 {
		return nil, 404, errors.New("тег на модерации не найден среди заявок под вашим рассмотрением")
	}

	return &result, 0, nil
}

func createTag(db *gorm.DB, tag models.Tag) error {
	return db.Create(&tag).Error
}