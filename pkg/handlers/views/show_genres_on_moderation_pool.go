package views

import "github.com/gin-gonic/gin"

func (h handler) ShowGenresOnModerationPool(c *gin.Context) {
	c.HTML(200, "genres_on_moderation_pool.html", nil)
}
