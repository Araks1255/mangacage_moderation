package authors

import (
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"gorm.io/gorm"
)

func DeleteAuthorOnModeration(db *gorm.DB, authorOnModerationID, moderatorID uint) (deleted *models.AuthorOnModeration, code int, err error) {
	var deletedAuthorOnModeration models.AuthorOnModeration

	err = db.Raw(
		"DELETE FROM authors_on_moderation WHERE id = ? AND moderator_id = ? RETURNING name, creator_id",
		authorOnModerationID, moderatorID,
	).Scan(&deletedAuthorOnModeration).Error

	if err != nil {
		return nil, 500, err
	}

	if deletedAuthorOnModeration.Name == "" {
		return nil, 404, errors.New("автор на модерации не найден среди рассматриваемых вами")
	}

	return &deletedAuthorOnModeration, 0, nil
}
