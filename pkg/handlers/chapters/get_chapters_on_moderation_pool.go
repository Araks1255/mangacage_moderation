package chapters

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	moderationDTO "github.com/Araks1255/mangacage_moderation/pkg/common/dto"
	"github.com/gin-gonic/gin"
)

type getChaptersPoolParams struct {
	dto.CommonParams

	NumberOfPagesFrom *int `form:"numberOfPagesFrom"`
	NumberOfPagesTo   *int `form:"numberOfPagesTo"`

	Volume *uint `form:"volume"`

	TitleID             *uint `form:"titleId"`
	TitleOnModerationID *uint `form:"titleOnModerationId"`
	CreatorID           *uint

	ModerationType string `form:"type"`
}

func (h handler) GetChaptersOnModerationPool(c *gin.Context) {
	var params getChaptersPoolParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Page - 1) * int(params.Limit)
	if offset < 0 {
		offset = 0
	}

	query := h.DB.Table("chapters_on_moderation AS com").
		Select("com.id, com.name, com.title_id, com.title_on_moderation_id, t.name AS title, tom.name AS title_on_moderation").
		Joins("LEFT JOIN titles AS t ON com.title_id = t.id").
		Joins("LEFT JOIN titles_on_moderation AS tom ON com.title_on_moderation_id = tom.id").
		Where("com.moderator_id IS NULL").
		Offset(offset).
		Limit(int(params.Limit))

	if params.Query != nil {
		query = query.Where("lower(com.name) ILIKE lower(?)", fmt.Sprintf("%%%s%%", *params.Query))
	}

	if params.CreatorID != nil {
		query = query.Where("com.creator_id = ?", params.CreatorID)
	}

	if params.ModerationType == "new" {
		query = query.Where("com.existing_id IS NULL")
	}
	if params.ModerationType == "edited" {
		query = query.Where("com.existing_id IS NOT NULL")
	}

	if params.NumberOfPagesFrom != nil {
		query = query.Where("com.number_of_pages >= ?", params.NumberOfPagesFrom)
	}
	if params.NumberOfPagesTo != nil {
		query = query.Where("com.number_of_pages <= ?", params.NumberOfPagesTo)
	}

	if params.Volume != nil {
		query = query.Where("com.volume = ?", params.Volume)
	}

	if params.TitleID != nil {
		query = query.Where("com.title_id = ?", params.TitleID)
	}
	if params.TitleOnModerationID != nil {
		query = query.Where("com.title_on_moderation_id = ?", params.TitleOnModerationID)
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	switch params.Sort {
	case "createdAt":
		query = query.Order(fmt.Sprintf("com.id %s", params.Order))
	case "numberOfPages":
		query = query.Order(fmt.Sprintf("com.number_of_pages %s", params.Order))
	default:
		query = query.Order(fmt.Sprintf("com.name %s", params.Order))
	}

	var result []moderationDTO.ChapterOnModerationDTO

	if err := query.Scan(&result).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "по вашему запросу ничего не найдено"})
		return
	}

	c.JSON(200, &result)
}
