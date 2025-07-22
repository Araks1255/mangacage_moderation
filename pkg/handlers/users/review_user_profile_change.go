package users

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) ReviewUserProfileChange(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	userProfileChangesOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id заявки на изменение профиля"})
		return
	}

	code, err := reviewUserProfileChanges(h.DB, uint(userProfileChangesOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на изменение профиля успешно взята вами на рассмотрение"})
}

func reviewUserProfileChanges(db *gorm.DB, userProfileChangesOnModerationID, moderatorID uint) (code int, err error) {
	result := db.Exec(
		"UPDATE users_on_moderation SET moderator_id = ? WHERE id = ? AND moderator_id IS NULL",
		moderatorID, userProfileChangesOnModerationID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("заявка не найдена в пуле изменений профиля")
	}

	return 0, nil
}
