package handlers

import (
	"net/http"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"

	"github.com/gin-gonic/gin"
)

type ServiceStatus string

const (
	UP   ServiceStatus = "ok"
	DOWN ServiceStatus = "down"
)

type StatusResponse struct {
	Status      ServiceStatus
	ServiceName string
	UpTime      string
	Environment string
	Version     string
}

type StatusController struct {
	dbMgr db.MongoManager
	lgr   *logger.AppLogger
}

func NewStatusController(lgr *logger.AppLogger, m db.MongoManager) *StatusController {
	return &StatusController{
		dbMgr: m,
		lgr:   lgr,
	}
}

// CheckStatus - Checks the health of all the dependencies of the service to ensure complete serviceability.
func (s *StatusController) CheckStatus(c *gin.Context) {
	var stat ServiceStatus
	var code int

	if err := s.dbMgr.Ping(); err == nil {
		stat = UP
		code = http.StatusOK
	} else {
		s.lgr.Error().Msg("unable to ping DB")
		stat = DOWN
		code = http.StatusFailedDependency
	}

	// send response
	c.JSON(code, stat)
}
