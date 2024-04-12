package middleware

import "github.com/gin-gonic/gin"

func InternalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Generally we would valid JWT token
		// TODO: Check for additional permissions for internal apis
		c.Next()
	}
}
