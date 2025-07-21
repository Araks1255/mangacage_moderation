package authors

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) DeclineAuthor(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	authorOnModerationID, reason, err := parseDeclineAuthorBody(c.ShouldBindJSON, c.Param)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	code, err := deleteAuthorOnModeration(h.DB, authorOnModerationID, claims.ID)
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

func deleteAuthorOnModeration(db *gorm.DB, authorOnModerationID, userID uint) (code int, err error) {
	result := db.Exec(
		"DELETE FROM authors_on_moderation WHERE id = ? AND moderator_id = ?",
		authorOnModerationID, userID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("автор на модерации не найден среди рассматриваемых вами")
	}

	return 0, nil
}
