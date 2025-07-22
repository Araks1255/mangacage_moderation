package tags

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) ReviewTagOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	tagOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тега на модерации"})
		return
	}

	code, err := reviewTagOnModeration(h.DB, uint(tagOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "тег успешно взят вами на рассмотрение"})
}

func reviewTagOnModeration(db *gorm.DB, tagOnModerationID, moderatorID uint) (code int, err error) {
	result := db.Exec(
		"UPDATE tags_on_moderation SET moderator_id = ? WHERE id = ? AND moderator_id IS NULL",
		moderatorID, tagOnModerationID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("заявка не найдена в пуле тегов на модерации")
	}

	return 0, nil
}
