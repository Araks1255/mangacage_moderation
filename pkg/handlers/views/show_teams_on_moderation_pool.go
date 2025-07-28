package views

import "github.com/gin-gonic/gin"

func (h handler) ShowTeamsOnModerationPool(c *gin.Context) {
	c.HTML(200, "teams_on_moderation_pool.html", nil)
}
