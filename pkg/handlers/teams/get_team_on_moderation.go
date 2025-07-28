package teams

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage_moderation/pkg/common/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeamOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	teamOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id команды на модерации"})
		return
	}

	var result dto.TeamOnModerationDTO

	err = h.DB.Table("teams_on_moderation AS tom").
		Select("tom.*, t.name AS existing, u.user_name AS creator").
		Joins("LEFT JOIN teams AS t ON tom.existing_id = t.id").
		Joins("INNER JOIN users AS u ON u.id = tom.creator_id").
		Where("tom.id = ?", teamOnModerationID).
		Where("tom.moderator_id = ?", claims.ID).
		Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда на модерации не найдена среди рассматриваемых вами"})
		return
	}

	c.JSON(200, &result)
}
