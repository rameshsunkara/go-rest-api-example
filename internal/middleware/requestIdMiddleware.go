package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIdentifier = "x-trace-id"

func ReqIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := c.Request.Header.Get(RequestIdentifier)
		if reqId == "" {
			reqId = uuid.New().String()
		}
		ctx := context.WithValue(c.Request.Context(), RequestIdentifier, reqId)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(RequestIdentifier, reqId)
		c.Next()
	}
}
