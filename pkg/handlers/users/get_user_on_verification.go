package users

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetUserOnVerification(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	userOnVerificationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id пользователя на верификации"})
		return
	}

	var result dto.ResponseUserDTO

	err = h.DB.Raw("SELECT * from users WHERE NOT verificated AND id = ? AND moderator_id = ?", userOnVerificationID, claims.ID).Scan(&result).Error
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "пользователь на верификации не найден среди рассматриваемых вами"})
		return
	}

	c.JSON(200, &result)
}
