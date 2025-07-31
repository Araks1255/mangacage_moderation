package feedback

import (
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
)

type handler struct {
	NotificationsClient pb.ModerationNotificationsClient
}

func RegisterRoutes(notificationsClient pb.ModerationNotificationsClient, r *gin.Engine) {
	h := handler{NotificationsClient: notificationsClient}

	r.POST("/moderation/api/feedback/users/:id", h.SendMessageToUser)
}
