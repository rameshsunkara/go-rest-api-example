package middleware

import "github.com/gin-gonic/gin"

func InternalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
