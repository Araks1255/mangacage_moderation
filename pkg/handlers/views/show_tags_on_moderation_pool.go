package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTagsOnModerationPool(c *gin.Context) {
	c.HTML(200, "tags_on_moderation_pool.html", nil)
}
