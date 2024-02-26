package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
)

func RequestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rLogger, _ := logger.ReqLogger(c)
		start := time.Now()
		c.Next()
		end := time.Now()
		rLogger.Info().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.String()).
			Str("path", c.Request.URL.Path).
			Int("responseStatus", c.Writer.Status()).
			Dur("responseInMS", end.Sub(start)).
			Interface("reqHeaders", c.Request.Header).
			Send()
	}
}
