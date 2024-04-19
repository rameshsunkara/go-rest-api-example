package server

import (
	"io"
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/rameshsunkara/go-rest-api-example/internal/log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/rameshsunkara/go-rest-api-example/internal/middleware"
	"github.com/rameshsunkara/go-rest-api-example/internal/types"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
)

var _runOnce sync.Once

func StartService(svcEnv types.ServiceEnv, dbMgr db.MongoManager, lgr *log.Logger) {
	_runOnce.Do(func() {
		r := WebRouter(svcEnv, dbMgr, lgr)
		err := r.Run(":" + svcEnv.Port)
		if err != nil {
			panic(err)
		}
	})
}

func WebRouter(svcEnv types.ServiceEnv, dbMgr db.MongoManager, lgr *log.Logger) *gin.Engine {
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

	internalAPIGrp := router.Group("/internal")
	internalAPIGrp.Use(middleware.AuthMiddleware())
	pprof.RouteRegister(internalAPIGrp, "pprof")
	router.GET("/metrics", gin.WrapH(promhttp.Handler())) // /metrics
	status := handlers.NewStatusController(dbMgr)
	router.GET("/status", status.CheckStatus) // /status

	// Dependencies for handlers
	d := dbMgr.Database()
	orders := db.NewOrdersRepo(d)

	// This is a dev mode only route to seed the local db
	if util.IsDevMode(svcEnv.Name) {
		seed := handlers.NewSeedController(orders)
		internalAPIGrp.POST("/seed-local-db", seed.SeedDB) // /seedDB
	}

	// Routes - Ecommerce
	externalAPIGrp := router.Group("/ecommerce/v1")
	externalAPIGrp.Use(middleware.AuthMiddleware())
	{
		ordersGroup := externalAPIGrp.Group("orders")
		{
			orders := handlers.NewOrdersController(orders)
			ordersGroup.GET("", orders.GetAll)
			ordersGroup.GET("/:id", orders.GetById)
			ordersGroup.POST("", orders.Post)
			ordersGroup.PUT("", orders.Post)
			ordersGroup.DELETE("/:id", orders.DeleteById)
		}
	}

	lgr.ZLogger.Info().Msg("Registered routes")
	for _, item := range router.Routes() {
		lgr.ZLogger.Info().
			Str("method", item.Method).
			Str("path", item.Path).
			Send()
	}
	return router
}
