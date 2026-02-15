package middleware

import (
	"time"

	"github.com/bogdanutanu/go-rest-api-example/pkg/flightrecorder"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/gin-gonic/gin"
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
