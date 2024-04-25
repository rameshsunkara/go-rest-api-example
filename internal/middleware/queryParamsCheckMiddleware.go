package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/external"
)

var GetOrdersListReqParams = map[string]bool{
	"limit":  true,
	"offset":   true,
}

var AllowedQueryParams = map[string]map[string]bool{
	http.MethodGet + "/ecommerce/v1/orders": GetOrdersListReqParams,
	http.MethodPost + "/ecommerce/v1/orders": nil,
	http.MethodGet + "/ecommerce/v1/orders/:id": nil,
	http.MethodDelete + "/ecommerce/v1/orders/:id": nil,
}

// QueryParamsCheckMiddleware - Middleware to check for unsupported query parameters
func QueryParamsCheckMiddleware(lgr *logger.AppLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, requestID := lgr.WithReqID(c)
		// Validate query params
		hasBadReqParams := HasUnSupportedQueryParams(c.Request, AllowedQueryParams[c.Request.Method+c.FullPath()])
		if hasBadReqParams {
			l.Info().Str("url path", c.Request.URL.Path).Msg("QueryParamsCheckMiddleware")
			apiErr := &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "",
				Message:        "Invalid query params",
				DebugID:        requestID,
			}
			l.Error().
				Int("HttpStatusCode", apiErr.HTTPStatusCode).
				Str("ErrorCode", apiErr.ErrorCode).
				Msg(apiErr.Message)
			c.AbortWithStatusJSON(apiErr.HTTPStatusCode, apiErr)
			return
		}
		c.Next()
	}
}

func HasUnSupportedQueryParams(req *http.Request, supportedParams map[string]bool) bool {
	queryParams := req.URL.Query()
	// Check for unsupported parameters
	for param := range queryParams {
		if _, ok := supportedParams[param]; !ok {
			// Handle the case of an unsupported parameter
			return true
		}
	}
	return false
}

