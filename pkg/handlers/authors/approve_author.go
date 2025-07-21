package authors

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

func (h handler) ApproveAuthor(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	authorOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id автора на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	authorOnModeration, code, err := popAuthorOnModeration(tx, uint(authorOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := createAuthor(tx, authorOnModeration.ToAuthor()); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "заявка на модерацию автора успешно одобрена"})
	// Уведомление
}

func popAuthorOnModeration(db *gorm.DB, authorOnModerationID, userID uint) (author *models.AuthorOnModeration, code int, err error) {
	var result models.AuthorOnModeration

	err = db.Raw(
		"DELETE FROM authors_on_moderation WHERE id = ? AND moderator_id = ?",
		authorOnModerationID, userID,
	).Scan(&result).Error

	if err != nil {
		return nil, 500, err
	}

	if result.ID == 0 {
		return nil, 404, errors.New("автор на модерации не найден среди заявок под вашим рассмотрением")
	}

	return &result, 0, nil
}

func createAuthor(db *gorm.DB, author models.Author) error {
	return db.Create(&author).Error
}
