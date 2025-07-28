package views

import (
	"github.com/gin-gonic/gin"
)

func (h handler) ShowUsersProfileChangesReviewingByMeCatalog(c *gin.Context) {
	c.HTML(200, "users_profile_changes_reviewing_by_me_catalog.html", nil)
}
