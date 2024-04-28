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

func StartService(svcEnv models.ServiceEnv, dbMgr db.MongoManager, lgr *logger.AppLogger) {
	startOnce.Do(func() {
		r := WebRouter(svcEnv, dbMgr, lgr)
		err := r.Run(":" + svcEnv.Port)
		if err != nil {
			panic(err)
		}
	})
}

func WebRouter(svcEnv models.ServiceEnv, dbMgr db.MongoManager, lgr *logger.AppLogger) *gin.Engine {
	ginMode := gin.ReleaseMode
	if util.IsDevMode(svcEnv.Name) {
		ginMode = gin.DebugMode
		gin.ForceConsoleColor()
	}
	gin.SetMode(ginMode)
	gin.EnableJsonDecoderDisallowUnknownFields()

	// Middleware
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.ReqIDMiddleware())
	router.Use(middleware.ResponseHeadersMiddleware())
	router.Use(middleware.RequestLogMiddleware(lgr))
	router.Use(gin.Recovery())

	internalAPIGrp := router.Group("/internal")
	internalAPIGrp.Use(middleware.AuthMiddleware())
	pprof.RouteRegister(internalAPIGrp, "pprof")
	router.GET("/metrics", gin.WrapH(promhttp.Handler())) // /metrics
	status := handlers.NewStatusController(dbMgr)
	router.GET("/status", status.CheckStatus) // /status

	// Dependencies for handlers
	d := dbMgr.Database()
	orders := db.NewOrdersRepo(d, lgr)

	// This is a dev mode only route to seed the local db
	if util.IsDevMode(svcEnv.Name) {
		seed := handlers.NewDataSeedHandler(orders)
		internalAPIGrp.POST("/seed-local-db", seed.SeedDB) // /seedDB
	}

	// Routes - Ecommerce
	externalAPIGrp := router.Group("/ecommerce/v1")
	externalAPIGrp.Use(middleware.AuthMiddleware())
	externalAPIGrp.Use(middleware.QueryParamsCheckMiddleware(lgr))
	{
		ordersGroup := externalAPIGrp.Group("orders")
		{
			orders := handlers.NewOrdersHandler(orders, lgr)
			ordersGroup.GET("", orders.GetAll)
			ordersGroup.GET("/:id", orders.GetByID)
			ordersGroup.POST("", orders.Create)
			ordersGroup.DELETE("/:id", orders.DeleteByID)
		}
	}

	lgr.Info().Msg("Registered routes")
	for _, item := range router.Routes() {
		lgr.Info().
			Str("method", item.Method).
			Str("path", item.Path).
			Send()
	}
	return router
}
