package server

import (
	"io"
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/rameshsunkara/go-rest-api-example/internal/middleware"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
)

var startOnce sync.Once

func Start(svcEnv *models.ServiceEnv, lgr *logger.AppLogger, dbMgr db.MongoManager) error {
	var err error
	var r *gin.Engine
	startOnce.Do(func() {
		r, err = WebRouter(svcEnv, lgr, dbMgr)
		if err != nil {
			return
		}
		err = r.Run(":" + svcEnv.Port)
	})
	return err
}

func WebRouter(svcEnv *models.ServiceEnv, lgr *logger.AppLogger, dbMgr db.MongoManager) (*gin.Engine, error) {
	ginMode := gin.ReleaseMode
	if util.IsDevMode(svcEnv.Name) {
		ginMode = gin.DebugMode
		gin.ForceConsoleColor()
	}
	gin.SetMode(ginMode)
	gin.EnableJsonDecoderDisallowUnknownFields()

	// Middleware
	gin.DefaultWriter = io.Discard
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.ReqIDMiddleware())
	router.Use(middleware.ResponseHeadersMiddleware())
	router.Use(middleware.RequestLogMiddleware(lgr))

	internalAPIGrp := router.Group("/internal")
	internalAPIGrp.Use(middleware.InternalAuthMiddleware()) // use special auth middleware to handle internal employees
	pprof.RouteRegister(internalAPIGrp, "pprof")
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	status, sHandlerErr := handlers.NewStatusHandler(lgr, dbMgr)
	if sHandlerErr != nil {
		return nil, sHandlerErr
	}
	router.GET("/status", status.CheckStatus)

	d := dbMgr.Database()
	ordersRepo, ordersRepoErr := db.NewOrdersRepo(lgr, d)
	if ordersRepoErr != nil {
		return nil, ordersRepoErr
	}

	// This is a dev mode only endpoint (route) to seed the local db
	if util.IsDevMode(svcEnv.Name) {
		if seed, seedHandlerErr := handlers.NewDataSeedHandler(lgr, ordersRepo); seedHandlerErr != nil {
			lgr.Error().Err(seedHandlerErr).Msg("seed-local-db endpoint will not be available")
		} else {
			internalAPIGrp.POST("/seed-local-db", seed.SeedDB)
		}
	}

	// Routes - Ecommerce
	externalAPIGrp := router.Group("/ecommerce/v1")
	externalAPIGrp.Use(middleware.AuthMiddleware())
	externalAPIGrp.Use(middleware.QueryParamsCheckMiddleware(lgr))
	ordersGroup := externalAPIGrp.Group("orders")
	ordersHandler, ordersHandlerErr := handlers.NewOrdersHandler(lgr, ordersRepo)
	if ordersHandlerErr != nil {
		return nil, ordersHandlerErr
	}
	ordersGroup.GET("", ordersHandler.GetAll)
	ordersGroup.GET("/:id", ordersHandler.GetByID)
	ordersGroup.POST("", ordersHandler.Create)
	ordersGroup.DELETE("/:id", ordersHandler.DeleteByID)

	lgr.Info().Msg("Registered routes")
	for _, item := range router.Routes() {
		lgr.Info().
			Str("method", item.Method).
			Str("path", item.Path).
			Send()
	}
	return router, nil
}
