package teams

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) ReviewTeamOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	teamOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id команды на модерации"})
		return
	}

	code, err := reviewTeam(h.DB, uint(teamOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "команда успешно взята вами на рассмотрение"})
}

func reviewTeam(db *gorm.DB, teamOnModerationID, moderatorID uint) (code int, err error) {
	result := db.Exec(
		"UPDATE teams_on_moderation SET moderator_id = ? WHERE id = ? AND moderator_id IS NULL",
		moderatorID, teamOnModerationID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("заявка не найдена в пуле команд на модерации")
	}

	return 0, nil
}
