package views

import "github.com/gin-gonic/gin"

func (h handler) ShowChaptersOnModerationPool(c *gin.Context) {
	c.HTML(200, "chapters_on_moderation_pool.html", nil)
}