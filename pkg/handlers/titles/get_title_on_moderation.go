package titles

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage_moderation/pkg/common/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTitleOnModeraiton(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла на модерации"})
		return
	}

	var result dto.TitleOnModerationDTO

	err = h.DB.Table("titles_on_moderation AS tom").
		Select(
			`tom.*, a.name AS author, aom.name AS author_on_moderation, t.name AS existing,
			ARRAY_AGG(DISTINCT g.name) AS genres, ARRAY_AGG(DISTINCT tags.name) AS tags`,
		).
		Joins("LEFT JOIN authors_on_moderation AS aom ON tom.author_on_moderation_id = aom.id").
		Joins("LEFT JOIN authors AS a ON tom.author_id = a.id").
		Joins("INNER JOIN title_on_moderation_genres AS tomg ON tomg.title_on_moderation_id = tom.id").
		Joins("INNER JOIN genres AS g ON g.id = tomg.genre_id").
		Joins("INNER JOIN title_on_moderation_tags AS tomt ON tomt.title_on_moderation_id = tom.id").
		Joins("INNER JOIN tags ON tags.id = tomt.tag_id").
		Where("tom.id = ?", titleOnModerationID).
		Where("tom.moderator_id = ?", claims.ID).
		Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл на модерации не найден среди рассматриваемых вами"})
		return
	}

	c.JSON(200, &result)
}
