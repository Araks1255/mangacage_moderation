package views

import "github.com/gin-gonic/gin"

func (h handler) ShowUsersOnVerificationPool(c *gin.Context) {
	c.HTML(200, "users_on_verification_pool.html", nil)
}
