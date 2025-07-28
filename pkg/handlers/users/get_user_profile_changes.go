package users

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetUserProfileChanges(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	userOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id изменений профиля"})
		return
	}

	var result dto.ResponseUserDTO

	err = h.DB.Raw(
		`SELECT
			uom.*, u.user_name AS existing
		from
			users_on_moderation AS uom
			INNER JOIN users AS u ON u.id = uom.existing_id
		WHERE
			uom.id = ? AND uom.moderator_id = ?`,
		userOnModerationID, claims.ID,
	).Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "изменения профиля не найдены среди рассматриваемых вами"})
		return
	}

	c.JSON(200, &result)
}
