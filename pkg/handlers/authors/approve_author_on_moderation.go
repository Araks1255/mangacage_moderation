package authors

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage_moderation/pkg/handlers/helpers/authors"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) ApproveAuthorOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	authorOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id автора на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	authorID, creatorID, code, err := createAuthorFromAuthorOnModerationByID(tx, uint(authorOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	err = replaceTitlesOnModerationAuthorOnModerationID(tx, uint(authorOnModerationID), authorID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if _, _, err := authors.DeleteAuthorOnModeration(tx, uint(authorOnModerationID), claims.ID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "заявка на модерацию автора успешно одобрена"})

	if _, err := h.NotificationsClient.NotifyAboutApprovedModerationRequest(
		c.Request.Context(), &pb.ApprovedEntity{
			Entity:    enums.Entity_ENTITY_AUTHOR,
			ID:        uint64(authorID),
			CreatorID: uint64(creatorID),
		},
	); err != nil {
		log.Println(err)
	}
}

func createAuthorFromAuthorOnModerationByID(db *gorm.DB, authorOnModerationID, moderatorID uint) (authorID, creatorID uint, code int, err error) {
	var authorOnModeration models.AuthorOnModeration

	err = db.Raw(
		"SELECT * FROM authors_on_moderation WHERE id = ? AND moderator_id = ?",
		authorOnModerationID, moderatorID,
	).Scan(&authorOnModeration).Error

	if err != nil {
		return 0, 0, 500, err
	}

	if authorOnModeration.ID == 0 {
		return 0, 0, 404, errors.New("автор на модерации не найден среди рассматриваемых вами")
	}

	newAuthor := authorOnModeration.ToAuthor()

	if err := db.Create(&newAuthor).Error; err != nil {
		return 0, 0, 500, err
	}

	return newAuthor.ID, authorOnModeration.CreatorID, 0, nil
}

func replaceTitlesOnModerationAuthorOnModerationID(db *gorm.DB, authorOnModerationID, authorID uint) error {
	return db.Exec(
		`UPDATE titles_on_moderation SET
			author_id = ?,
			author_on_moderation_id = NULL
		WHERE
			author_on_moderation_id = ?`,
		authorID, authorOnModerationID,
	).Error
}
