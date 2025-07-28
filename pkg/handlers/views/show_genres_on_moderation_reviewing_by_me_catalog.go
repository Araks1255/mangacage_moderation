package views

import (
	"github.com/gin-gonic/gin"
)

func (h handler) ShowGenresOnModerationReviewingByMeCatalog(c *gin.Context) {
	c.HTML(200, "genres_on_moderation_reviewing_by_me_catalog.html", nil)
}
