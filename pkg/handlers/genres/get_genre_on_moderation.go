package genres

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage_moderation/pkg/common/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetGenreOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	genreOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id жанра на модерации"})
		return
	}

	var result dto.GenreOnModerationDTO

	err = h.DB.Table("genres_on_moderation AS gom").
		Select("gom.*, u.user_name AS creator").
		Joins("INNER JOIN users AS u ON u.id = gom.creator_id").
		Where("gom.id = ?", genreOnModerationID).
		Where("gom.moderator_id = ?", claims.ID).
		Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "жанр на модерации не найден среди рассматриваемых вами"})
		return
	}

	c.JSON(200, &result)
}

