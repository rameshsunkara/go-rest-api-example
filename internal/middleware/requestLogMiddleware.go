package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
)

func RequestLogMiddleware(lgr *logger.AppLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, _ := lgr.WithReqID(c)
		start := time.Now()
		c.Next()
		// consider adding more request identifiers such as userId etc.,
		l.Info().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.String()).
			Str("path", c.FullPath()).
			Str("userAgent", c.Request.UserAgent()).
			Int("respStatus", c.Writer.Status()).
			Dur("elapsedMs", time.Since(start)).
			Send()
	}
}
