package genres

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

func (h handler) ApproveGenreOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	genreOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id жанра на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	genreOnModeration, code, err := popGenreOnModeration(tx, uint(genreOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := createGenre(tx, genreOnModeration.ToGenre()); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "заявка на модерацию жанра успешно одобрена"})
	// Уведомление
}

func popGenreOnModeration(db *gorm.DB, genreOnModerationID, userID uint) (genre *models.GenreOnModeration, code int, err error) {
	var result models.GenreOnModeration

	err = db.Raw(
		"DELETE FROM genres_on_moderation WHERE id = ? AND moderator_id = ?",
		genreOnModerationID, userID,
	).Scan(&result).Error

	if err != nil {
		return nil, 500, err
	}

	if result.ID == 0 {
		return nil, 404, errors.New("жанр на модерации не найден среди заявок под вашим рассмотрением")
	}

	return &result, 0, nil
}

func createGenre(db *gorm.DB, genre models.Genre) error {
	return db.Create(&genre).Error
}
