package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIdentifier = "X-Request-ID"

func ReqIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.Request.Header.Get(RequestIdentifier)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(c.Request.Context(), RequestIdentifier, reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(RequestIdentifier, reqID)
		c.Next()
	}
}
