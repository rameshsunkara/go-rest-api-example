package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextKey is a type for context keys.
type ContextKey string

const (
	// RequestIdentifier is the header name for request ID.
	RequestIdentifier = "X-Request-ID"
)

// ReqIDMiddleware injects a request ID into the context and response header, creates one if it is not present already.
func ReqIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.Request.Header.Get(RequestIdentifier)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(c.Request.Context(), ContextKey(RequestIdentifier), reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(RequestIdentifier, reqID)
		c.Next()
	}
}
