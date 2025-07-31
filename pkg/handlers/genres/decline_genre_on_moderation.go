package genres

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) DeclineGenreOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	genreOnModerationID, reason, err := parseDeclineGenreOnModerationBody(c.ShouldBindJSON, c.Param)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	genreOnModeration, code, err := deleteGenreOnModeration(h.DB, genreOnModerationID, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на модерацию жанра успешно отклонена"})

	if _, err := h.NotificationsClient.SendModerationRequestDeclineReason(
		c.Request.Context(), &pb.ModerationRequestDeclineReason{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_GENRE,
			EntityName:         genreOnModeration.Name,
			CreatorID:          uint64(genreOnModeration.CreatorID),
			Reason:             reason,
		},
	); err != nil {
		log.Println(err)
	}
}

func parseDeclineGenreOnModerationBody(bindFn func(any) error, paramFn func(string) string) (genreID uint, reason string, err error) {
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

func deleteGenreOnModeration(db *gorm.DB, genreOnModerationID, userID uint) (deleted *models.GenreOnModeration, code int, err error) {
	var deletedGenreOnModeration models.GenreOnModeration

	err = db.Raw(
		"DELETE FROM genres_on_moderation WHERE id = ? AND moderator_id = ? RETURNING name, creator_id",
		genreOnModerationID, userID,
	).Scan(&deletedGenreOnModeration).Error

	if err != nil {
		return nil, 500, err
	}

	if deletedGenreOnModeration.Name == "" {
		return nil, 404, errors.New("жанр на модерации не найден среди рассматриваемых вами")
	}

	return &deletedGenreOnModeration, 0, nil
}
