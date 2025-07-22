package genres

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) ReviewGenreOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	genreOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id жанра на модерации"})
		return
	}

	code, err := reviewGenreOnModeration(h.DB, uint(genreOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "жанр успешно взят вами на рассмотрение"})
}

func reviewGenreOnModeration(db *gorm.DB, genreOnModerationID, moderatorID uint) (code int, err error) {
	result := db.Exec(
		"UPDATE genres_on_moderation SET moderator_id = ? WHERE id = ? AND moderator_id IS NULL",
		moderatorID, genreOnModerationID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("заявка не найдена в пуле жанров на модерации")
	}

	return 0, nil
}
