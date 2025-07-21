package users

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetUsersProfileChangesPool(c *gin.Context) {
	var params dto.CommonParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Page - 1) * int(params.Limit)
	if offset < 0 {
		offset = 0
	}

	query := h.DB.Table("users_on_moderation").
		Select("*").
		Where("moderator_id IS NULL").
		Offset(offset).
		Limit(int(params.Limit))

	if params.Query != nil {
		query = query.Where("lower(user_name) ILIKE lower(?)", fmt.Sprintf("%%%s%%", *params.Query))
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	switch params.Sort {
	case "createdAt":
		query = query.Order(fmt.Sprintf("id %s", params.Order))
	default:
		query = query.Order(fmt.Sprintf("name %s", *params.Query))
	}

	var result []dto.ResponseUserDTO

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
