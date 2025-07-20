package users

import "github.com/gin-gonic/gin"

func (h handler) GetUsersPool(c *gin.Context) {
	c.JSON(200, gin.H{"message": "GetUsersPool handler - not implemented"})
}
