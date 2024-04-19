package handlers

import (
	"net/http"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"

	"github.com/rs/zerolog/log"

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
}

func NewStatusController(m db.MongoManager) *StatusController {
	return &StatusController{
		dbMgr: m,
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
		log.Error().Msg("unable to connect to DB")
		stat = DOWN
		code = http.StatusFailedDependency
	}

	// send response
	c.JSON(code, stat)
}
