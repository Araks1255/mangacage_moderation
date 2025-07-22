package users

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) ReviewUserVerification(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	userOnVerificationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id заявки на верификацию аккаунта"})
		return
	}

	code, err := reviewUserVerification(h.DB, uint(userOnVerificationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на верификацию аккаунта успешно взята вами на рассмотрение"})
}

func reviewUserVerification(db *gorm.DB, userOnVerificationID, moderatorID uint) (code int, err error) {
	result := db.Exec(
		"UPDATE users SET moderator_id = ? WHERE id = ? AND moderator_id IS NULL AND NOT verificated",
		moderatorID, userOnVerificationID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("заявка не найдена в пуле аккаунтов на верификации")
	}

	return 0, nil
}
