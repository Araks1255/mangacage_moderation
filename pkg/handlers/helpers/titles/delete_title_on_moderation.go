package titles

import (
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"gorm.io/gorm"
)

func DeleteTitleOnModeration(db *gorm.DB, titleOnModerationID, moderatorID uint) (deleted *models.TitleOnModeration, code int, err error) {
	var deletedTitleOnModeration models.TitleOnModeration

	err = db.Raw(
		"DELETE FROM titles_on_moderation WHERE id = ? AND moderator_id = ? RETURNING id, name, creator_id",
		titleOnModerationID, moderatorID,
	).Scan(&deletedTitleOnModeration).Error

	if err != nil {
		return nil, 500, err
	}

	if deletedTitleOnModeration.ID == 0 {
		return nil, 404, errors.New("тайтл на модерации не найден среди рассматриваемых вами")
	}

	return &deletedTitleOnModeration, 0, nil
}
