package authors

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) ReviewAuthorOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	authorOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id автора на модерации"})
		return
	}

	code, err := reviewAuthorOnModeration(h.DB, uint(authorOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "автор успешно взят вами на рассмотрение"})
}

func reviewAuthorOnModeration(db *gorm.DB, authorOnModerationID, moderatorID uint) (code int, err error) {
	result := db.Exec(
		"UPDATE authors_on_moderation SET moderator_id = ? WHERE id = ? AND moderator_id IS NULL",
		moderatorID, authorOnModerationID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("заявка не найдена в пуле авторов на модерации")
	}

	return 0, nil
}
