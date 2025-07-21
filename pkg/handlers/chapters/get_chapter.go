package chapters

import (
	"strconv"

	"github.com/Araks1255/mangacage_moderation/pkg/common/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetChapter(c *gin.Context) {
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
		Scan(&result).Error

}
