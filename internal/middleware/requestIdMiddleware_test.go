package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/middleware"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestReqIDMiddleware(t *testing.T) {
	type reqIDMiddlewareTestCase struct {
		Description  string
		InputReqID   string
		InputReqPath string
	}

	var testCases = []reqIDMiddlewareTestCase{
		{
			Description:  "ensure request id is set when not provided",
			InputReqID:   "",
			InputReqPath: "/test/1",
		},
		{
			Description:  "ensure request id is set when provided",
			InputReqID:   "123",
			InputReqPath: "/test/2",
		},
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(resp)
	r.Use(middleware.ReqIDMiddleware())

	for _, tc := range testCases {
		var hasCorrectReqID bool
		r.GET(tc.InputReqPath, func(ctx *gin.Context) {
			if rID := ctx.Request.Context().Value(util.ContextKey(util.RequestIdentifier)); rID != nil {
				if rIDStr, ok := rID.(string); ok {
					reqIDPassed := len(tc.InputReqID) > 0
					if reqIDPassed && rID == tc.InputReqID || (!reqIDPassed && len(rIDStr) > 0) {
						hasCorrectReqID = true
					}
				}
			}
			ctx.String(200, "OK")
		})

		c.Request, _ = http.NewRequest(http.MethodGet, tc.InputReqPath, nil)
		c.Request.Header.Set(util.RequestIdentifier, tc.InputReqID)
		r.ServeHTTP(resp, c.Request)
		// Check response header
		assert.NotEmpty(t, resp.Header().Get(util.RequestIdentifier), tc.Description)
		// Check request id is set in context
		assert.True(t, hasCorrectReqID, tc.Description)
	}
}
