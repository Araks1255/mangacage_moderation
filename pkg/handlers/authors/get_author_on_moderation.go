package authors

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage_moderation/pkg/common/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetAuthorOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	authorOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id автора на модерации"})
		return
	}

	var result dto.AuthorOnModerationDTO

	err = h.DB.Raw(
		`SELECT
			aom.*, u.user_name AS creator
		FROM
			authors_on_moderation AS aom
			INNER JOIN users AS u ON u.id = aom.creator_id
		WHERE
			aom.id = ? AND aom.moderator_id = ?`,
		authorOnModerationID, claims.ID,
	).Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "автор на модерации не найден среди рассматриваемых вами"})
		return
	}

	c.JSON(200, &result)
}
