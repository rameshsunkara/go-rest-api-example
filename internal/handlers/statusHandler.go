package handlers

import (
	"errors"
	"net/http"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"

	"github.com/gin-gonic/gin"
)

type StatusHandler struct {
	dbMgr db.MongoManager
	lgr   *logger.AppLogger
}

func NewStatusHandler(lgr *logger.AppLogger, m db.MongoManager) (*StatusHandler, error) {
	if lgr == nil || m == nil {
		return nil, errors.New("missing required inputs to create status handler")
	}
	return &StatusHandler{
		dbMgr: m,
		lgr:   lgr,
	}, nil
}

// CheckStatus checks the health of all service dependencies to ensure full serviceability.
func (s *StatusHandler) CheckStatus(c *gin.Context) {
	var code int

	if err := s.dbMgr.Ping(); err == nil {
		code = http.StatusNoContent
	} else {
		s.lgr.Error().Msg("failed to ping DB")
		code = http.StatusFailedDependency
	}

	// Check the status of any other dependencies you may have here

	// send response
	c.JSON(code, nil)
}
