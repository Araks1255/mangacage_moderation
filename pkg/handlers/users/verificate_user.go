package users

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
)

func (h handler) VerificateUser(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id пользователя"})
		return
	}

	result := h.DB.Exec("UPDATE users SET verificated = true WHERE id = ? AND moderator_id = ?", userID, claims.ID)

	if result.Error != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "пользователь не найден среди верифицируемых вами"})
		return
	}

	c.JSON(200, gin.H{"success": "пользователь успешно верифицирован"})

	if _, err := h.NotificationsClient.NotifyUserAboutVerificatedAccount(
		c.Request.Context(), &pb.VerificatedUser{ID: userID},
	); err != nil {
		log.Println(err)
	}
}
