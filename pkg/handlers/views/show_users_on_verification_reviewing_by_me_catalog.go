package views

import (
	"github.com/gin-gonic/gin"
)

func (h handler) ShowUsersOnVerificationReviewingByMeCatalog(c *gin.Context) {
	c.HTML(200, "users_on_verification_reviewing_by_me_catalog.html", nil)
}
