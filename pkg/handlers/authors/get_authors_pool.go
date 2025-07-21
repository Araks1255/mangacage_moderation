package authors

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	moderationDTO "github.com/Araks1255/mangacage_moderation/pkg/common/dto"
	"github.com/gin-gonic/gin"
)

type getAuthorsPoolParams struct {
	dto.CommonParams

	CreatorID *uint `form:"creatorId"`
}

func (h handler) GetAuthorsPool(c *gin.Context) {
	var params getAuthorsPoolParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Page - 1) * int(params.Limit)
	if offset < 0 {
		offset = 0
	}

	query := h.DB.Table("authors_on_moderation").
		Where("moderator_id IS NULL").
		Offset(offset).
		Limit(int(params.Limit))

	if params.Query != nil {
		query = query.Where("lower(name) ILIKE lower(?)", fmt.Sprintf("%%%s%%", *params.Query))
	}

	if params.CreatorID != nil {
		query = query.Where("creator_id = ?", params.CreatorID)
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	switch params.Sort {
	case "createdAt":
		query = query.Order(fmt.Sprintf("id %s", params.Order))
	default:
		query = query.Order(fmt.Sprintf("name %s", params.Order))
	}

	var result []moderationDTO.AuthorOnModerationDTO

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
