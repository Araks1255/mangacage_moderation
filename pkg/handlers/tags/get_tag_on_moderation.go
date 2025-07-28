package tags

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage_moderation/pkg/common/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTagOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	tagOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тега на модерации"})
		return
	}

	var result dto.TagOnModerationDTO

	err = h.DB.Table("tags_on_moderation AS tom").
		Select("tom.*, u.user_name AS creator").
		Joins("INNER JOIN users AS u ON u.id = tom.creator_id").
		Where("tom.id = ?", tagOnModerationID).
		Where("tom.moderator_id = ?", claims.ID).
		Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тег на модерации не найден среди рассматриваемых вами"})
		return
	}

	c.JSON(200, &result)
}
