package titles

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) ReviewTitleOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла на модерации"})
		return
	}

	code, err := reviewTitle(h.DB, uint(titleOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "тайтл успешно взят вами на рассмотрение"})
}

func reviewTitle(db *gorm.DB, titleOnModerationID, moderatorID uint) (code int, err error) {
	result := db.Exec(
		"UPDATE titles_on_moderation SET moderator_id = ? WHERE id = ? AND moderator_id IS NULL",
		moderatorID, titleOnModerationID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("заявка не найдена в пуле тайтлов на модерации")
	}

	return 0, nil
}
