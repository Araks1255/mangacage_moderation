package authors

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage_moderation/pkg/common/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetAuthorsOnModerationReviewingByMe(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var result []dto.AuthorOnModerationDTO

	err := h.DB.Raw("SELECT * FROM authors_on_moderation WHERE moderator_id = ? ORDER BY id ASC", claims.ID).Scan(&result).Error
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, &result)
}
