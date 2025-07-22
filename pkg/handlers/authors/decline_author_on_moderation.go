package authors

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/helpers/authors"
	"github.com/gin-gonic/gin"
)

func (h handler) DeclineAuthorOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	authorOnModerationID, reason, err := parseDeclineAuthorBody(c.ShouldBindJSON, c.Param)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	code, err := authors.DeleteAuthorOnModeration(h.DB, authorOnModerationID, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на модерацию автора успешно отклонена"})
	// Уведомление с причиной
	log.Println(reason)
}

func parseDeclineAuthorBody(bindFn func(any) error, paramFn func(string) string) (authorID uint, reason string, err error) {
	var requestBody struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := bindFn(&requestBody); err != nil {
		return 0, "", err
	}

	id, err := strconv.ParseUint(paramFn("id"), 10, 64)
	if err != nil {
		return 0, "", err
	}

	return uint(id), requestBody.Reason, nil
}
