package users

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetUserProfileChangesReviewingByMe(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var result []dto.ResponseUserDTO

	err := h.DB.Raw("SELECT * FROM users_on_moderation WHERE moderator_id = ? ORDER BY id ASC", claims.ID).Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, &result)
}
