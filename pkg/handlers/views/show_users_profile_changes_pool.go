package views

import "github.com/gin-gonic/gin"

func (h handler) ShowUsersProfileChangesPool(c *gin.Context) {
	c.HTML(200, "users_profile_changes_pool.html", nil)
}
