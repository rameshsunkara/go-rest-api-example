package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/pkg/flightrecorder"
	"github.com/rameshsunkara/go-rest-api-example/pkg/logger"
)

const (
	// SlowRequestThreshold defines when to capture flight recorder traces.
	SlowRequestThreshold = 500 * time.Millisecond
)

func RequestLogMiddleware(lgr logger.Logger, fr *flightrecorder.Recorder) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, _ := lgr.WithReqID(c)
		start := time.Now()

		c.Next()

		elapsed := time.Since(start)

		// Log the request
		l.Info().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.String()).
			Str("path", c.FullPath()).
			Str("userAgent", c.Request.UserAgent()).
			Int("respStatus", c.Writer.Status()).
			Dur("elapsedMs", elapsed).
			Send()

		// Capture trace for slow requests
		if fr != nil && elapsed > SlowRequestThreshold {
			fr.CaptureSlowRequest(l, c.Request.Method, c.FullPath(), elapsed)
		}
	}
}
