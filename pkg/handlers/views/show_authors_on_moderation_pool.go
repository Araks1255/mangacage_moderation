package views

import "github.com/gin-gonic/gin"

func (h handler) ShowAuthorsOnModerationPool(c *gin.Context) {
	c.HTML(200, "authors_on_moderation_pool.html", nil)
}
