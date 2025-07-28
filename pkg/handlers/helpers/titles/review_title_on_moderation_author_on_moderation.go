package titles

import "gorm.io/gorm"

func ReviewTitleOnModerationAuthorOnModeration(db *gorm.DB, titleOnModerationID, moderatorID uint) error {
	if titleOnModerationID == 0 {
		return nil
	}
	return db.Exec(
		`UPDATE
			authors_on_moderation
		SET
			moderator_id = ?
		WHERE
			id = (SELECT author_on_moderation_id FROM titles_on_moderation WHERE id = ?)
		AND
			moderator_id IS NULL`,
		moderatorID, titleOnModerationID,
	).Error
}
