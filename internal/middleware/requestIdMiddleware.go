package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
)

// ReqIDMiddleware injects a request ID into the context and response header, creates one if it is not present already.
func ReqIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.Request.Header.Get(util.RequestIdentifier)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(c.Request.Context(), util.ContextKey(util.RequestIdentifier), reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(util.RequestIdentifier, reqID)
		c.Next()
	}
}
