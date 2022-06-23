package controllers

import (
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"net/http"
	"time"

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
	UpTime      time.Time
	Environment string
	Version     string
}

type StatusHandler struct {
	svcInfo models.ServiceInfo
	dbMgr   db.DataManager
}

func NewStatusHandler(s models.ServiceInfo, m db.DataManager) *StatusHandler {
	return &StatusHandler{
		svcInfo: s,
		dbMgr:   m,
	}
}

func (s *StatusHandler) CheckStatus(c *gin.Context) {
	log.Debug().Msg("in CheckStatus")
	var stat ServiceStatus
	var code int
	err := s.dbMgr.Ping()
	if err == nil {
		stat = UP
		code = http.StatusOK
	} else {
		log.Error().Msg("unable to connect to DB")
		stat = DOWN
		code = http.StatusFailedDependency
	}
	status := &StatusResponse{
		Status:      stat,
		ServiceName: s.svcInfo.Name,
		UpTime:      s.svcInfo.UpTime,
		Environment: s.svcInfo.Environment,
		Version:     s.svcInfo.Version,
	}
	c.JSON(code, status)
}
