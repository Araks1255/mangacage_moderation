package authors

import (
	"errors"

	"gorm.io/gorm"
)

func DeleteAuthorOnModeration(db *gorm.DB, authorOnModerationID, moderatorID uint) (code int, err error) {
	result := db.Exec(
		"DELETE FROM authors_on_moderation WHERE id = ? AND moderator_id = ?",
		authorOnModerationID, moderatorID,
	)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("автор на модерации не найден среди рассматриваемых вами")
	}

	return 0, nil
}
