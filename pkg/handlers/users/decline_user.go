package users

import "github.com/gin-gonic/gin"

func (h handler) DeclineUser(c *gin.Context) {
	c.JSON(200, gin.H{"message": "DeclineUser handler - not implemented"})
}
