package chapters

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage_moderation/pkg/common/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetChapterOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapterOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы на модерации"})
		return
	}

	var result dto.ChapterOnModerationDTO

	err = h.DB.Table("chapters_on_moderation AS com").
		Select("com.*, t.name AS title, tom.name AS title_on_moderation, teams.name AS team, u.user_name AS creator").
		Joins("LEFT JOIN titles AS t ON com.title_id = t.id").
		Joins("LEFT JOIN titles_on_moderation AS tom ON com.title_on_moderation_id = tom.id").
		Joins("LEFT JOIN teams ON com.team_id = teams.id").
		Joins("INNER JOIN users AS u ON u.id = com.creator_id").
		Where("com.id = ?", chapterOnModerationID).
		Where("com.moderator_id = ?", claims.ID).
		Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава на модерации не найдена среди рассматриваемых вами"})
		return
	}

	c.JSON(200, &result)
}
