package titles

import (
	"errors"

	"gorm.io/gorm"
)

func DeleteTitleOnModeration(db *gorm.DB, titleOnModerationID, moderatorID uint) (code int, err error) {
	result := db.Exec(
		"DELETE FROM titles_on_moderation WHERE id = ? AND moderator_id = ?",
		titleOnModerationID, moderatorID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("тайтл на модерации не найден среди рассматриваемых вами")
	}

	return 0, nil
}
