package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ResponseHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		csp := []string{"default-src: 'self'", "font-src: 'fonts.googleapis.com'", "frame-src: 'none'"}
		c.Writer.Header().Set("Content-Security-Policy", strings.Join(csp, "; "))
		c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
		c.Writer.Header().Set("X-XSS-Protection", "1")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; preload")
		c.Next()
	}
}

