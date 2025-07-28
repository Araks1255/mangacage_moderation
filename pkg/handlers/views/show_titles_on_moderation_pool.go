package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) ShowTitlesOnModerationPool(c *gin.Context) {
	c.HTML(http.StatusOK, "titles_on_moderation_pool.html", nil)
}
