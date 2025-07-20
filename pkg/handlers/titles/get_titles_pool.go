package titles

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type getTitlesPoolParams struct {
	dto.CommonParams

	PublishingStatus *string `form:"publishingStatus"`
	Type             *string `form:"type"`
	ModerationType   string  `form:"moderationType"`

	YearFrom     *int `form:"yearFrom"`
	YearTo       *int `form:"yearTo"`
	AgeLimitFrom *int `form:"ageLimitFrom"`
	AgeLimitTo   *int `form:"ageLimitTo"`

	Genres []string `form:"genres"`
	Tags   []string `form:"tags"`

	AuthorID             *uint `form:"authorId"`
	AuthorOnModerationID *uint `form:"authorOnModerationId"`
	CreatorID            *uint
}

func (h handler) GetTitlesPool(c *gin.Context) {
	var params getTitlesPoolParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Page - 1) * int(params.Limit)
	if offset < 0 {
		offset = 0
	}

	query := h.DB.Table("titles_on_moderation AS tom").
		Select("tom.*, teams.name AS team").
		Joins("INNER JOIN teams ON teams.id = tom.team_id").
		Where("tom.moderator_id IS NULL").
		Offset(offset).
		Limit(int(params.Limit))

	if params.Query != nil {
		query = query.Where("lower(tom.name) ILIKE lower(?)", fmt.Sprintf("%%%s%%", *params.Query))
	}

	if params.CreatorID != nil {
		query = query.Where("tom.creator_id = ?", params.CreatorID)
	}

	if params.PublishingStatus != nil {
		query = query.Where("tom.publishing_status = ?", params.PublishingStatus)
	}

	if params.Type != nil {
		query = query.Where("tom.type = ?", params.Type)
	}

	if params.AuthorID != nil {
		query = query.Where("tom.author_id = ?", params.AuthorID)
	}
	if params.AuthorOnModerationID != nil {
		query = query.Where("tom.author_on_moderation_id = ?", params.AuthorOnModerationID)
	}

	if params.YearFrom != nil {
		query = query.Where("tom.year_of_release >= ?", params.YearFrom)
	}
	if params.YearTo != nil {
		query = query.Where("tom.year_of_release <= ?", params.YearTo)
	}

	if params.AgeLimitFrom != nil {
		query = query.Where("tom.age_limit >= ?", params.AgeLimitFrom)
	}
	if params.AgeLimitTo != nil {
		query = query.Where("tom.age_limit <= ?", params.AgeLimitTo)
	}

	switch params.ModerationType {
	case "new":
		query = query.Where("tom.existing_id IS NULL")
	case "edited":
		query = query.Where("tom.existing_id IS NOT NULL")
	}

	if len(params.Genres) > 0 {
		query = query.Where(`(
			SELECT ARRAY(
				SELECT g.name FROM genres AS g
				INNER JOIN title_on_moderation_genres AS tomg ON tomg.genre_id = g.id
				WHERE tomg.title_on_moderation_id = tom.id
			)
		)::TEXT[] @> ?::TEXT[]`, pq.Array(params.Genres))
	}

	if len(params.Tags) > 0 {
		query = query.Where(`(
			SELECT ARRAY(
				SELECT tags.name FROM tags
				INNER JOIN title_on_moderation_tags AS tomt ON tomt.tag_id = tags.id
				WHERE tomt.title_on_moderation_id = tom.id
			)
		)::TEXT[] @> ?::TEXT[]`, pq.Array(params.Tags))
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	switch params.Sort {
	case "createdAt":
		query = query.Order(fmt.Sprintf("tom.id %s", params.Order))
	default:
		query = query.Order(fmt.Sprintf("tom.name %s", params.Order))
	}

	var result []dto.ResponseTitleDTO

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
