package feedback

import (
	"strconv"

	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
)

func (h handler) SendMessageToUser(c *gin.Context) {
	receiverID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id получателя"})
		return
	}

	var requestBody struct {
		Message          string `json:"message" binding:"required"`
		AttachedEntity   string `json:"entity"`
		AttachedEntityID uint   `json:"entityId" binding:"required_with=AttachedEntity"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var attachedEntityEnumType enums.EntityOnModeration

	switch requestBody.AttachedEntity {
	case "title":
		attachedEntityEnumType = enums.EntityOnModeration_ENTITY_ON_MODERATION_TITLE
	case "chapter":
		attachedEntityEnumType = enums.EntityOnModeration_ENTITY_ON_MODERATION_CHAPTER
	}

	_, err = h.NotificationsClient.SendMessageToUser(
		c.Request.Context(), &pb.MessageFromModerator{
			ReceiverID:           receiverID,
			EntityOnModeration:   attachedEntityEnumType,
			EntityOnModerationID: uint64(requestBody.AttachedEntityID),
			Text:                 requestBody.Message,
		},
	)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "сообщение отправлено"})
}
