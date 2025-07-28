package chapters

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/helpers/titles"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) ReviewChapterOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapterOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	code, err := reviewChapterOnModeration(h.DB, uint(chapterOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "глава успешно взята вами на рассмотрение"})
}

func reviewChapterOnModeration(db *gorm.DB, chapterOnModerationID, moderatorID uint) (code int, err error) {
	result := db.Exec(
		"UPDATE chapters_on_moderation SET moderator_id = ? WHERE id = ? AND moderator_id IS NULL",
		moderatorID, chapterOnModerationID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("заявка не найдена в пуле глав на модерации")
	}

	chapterTitleOnModerationID, err := reviewChapterOnModerationTitleOnModeration(db, chapterOnModerationID, moderatorID)

	if err != nil {
		return 500, err
	}

	if err := titles.ReviewTitleOnModerationAuthorOnModeration(db, chapterTitleOnModerationID, moderatorID); err != nil {
		return 500, err
	}

	return 0, nil
}

func reviewChapterOnModerationTitleOnModeration(db *gorm.DB, chapterOnModerationID, moderatorID uint) (titleOnModerationID uint, err error) {
	var id *uint

	err = db.Raw(
		`UPDATE
			titles_on_moderation
		SET
			moderator_id = ?
		WHERE
			id = (SELECT title_on_moderation_id FROM chapters_on_moderation WHERE id = ?)
		AND
			moderator_id IS NULL
		RETURNING
			id`,
		moderatorID, chapterOnModerationID,
	).Scan(&id).Error

	if id != nil {
		titleOnModerationID = *id
	}

	return titleOnModerationID, err
}
