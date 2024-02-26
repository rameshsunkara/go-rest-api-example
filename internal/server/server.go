package server

import (
	"io"
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"

	// TODO: "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/rameshsunkara/go-rest-api-example/internal/middleware"
	"github.com/rameshsunkara/go-rest-api-example/internal/types"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
)

var RunOnce sync.Once

func StartService(svcInfo types.ServiceInfo, svcEnv types.ServiceEnv, dbMgr db.MongoManager, lgr *logger.Logger) {
	RunOnce.Do(func() {
		r := WebRouter(svcInfo, dbMgr, lgr)
		err := r.Run(":" + svcEnv.Port)
		if err != nil {
			panic(err)
		}
	})
}

func WebRouter(svcInfo types.ServiceInfo, dbMgr db.MongoManager, lgr *logger.Logger) *gin.Engine {
	ginMode := gin.ReleaseMode
	if util.IsDevMode(svcInfo.Environment) {
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

	// Routes

	// Routes - Status Check

	adminGroup := router.Group("/internal")
	adminGroup.Use()
	pprof.RouteRegister(adminGroup, "pprof")
	// TODO: router.GET("/metrics", gin.WrapH(promhttp.Handler())) // /metrics
	status := handlers.NewStatusController(svcInfo, dbMgr)
	router.GET("/status", status.CheckStatus) // /status

	// Dependencies for handlers
	d := dbMgr.Database()
	orders := db.NewOrdersRepo(d)

	// Routes - Seed DB
	if util.IsDevMode(svcInfo.Environment) {
		seed := handlers.NewSeedController(orders)
		router.POST("/seedDB", seed.SeedDB) // /seedDB
	}

	// Routes - Orders
	v1 := router.Group("/api/v1")
	{
		ordersGroup := v1.Group("orders")
		{
			orders := handlers.NewOrdersController(orders)
			ordersGroup.GET("", orders.GetAll)            // api/v1/orders
			ordersGroup.GET("/:id", orders.GetById)       // api/v1/orders/:id
			ordersGroup.POST("", orders.Post)             // api/v1/orders
			ordersGroup.PUT("", orders.Post)              // api/v1/orders
			ordersGroup.DELETE("/:id", orders.DeleteById) // api/v1/orders/:id
		}
	}
	return router
}
