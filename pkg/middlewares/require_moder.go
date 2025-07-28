package middlewares

import (
	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func RequireModer(modersIDs map[uint]struct{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(*auth.Claims)

		if _, ok := modersIDs[claims.ID]; !ok {
			c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для доступа к сервису модерации"})
			return
		}

		c.Next()
	}
}
