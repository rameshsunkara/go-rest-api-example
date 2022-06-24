package server

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/controllers"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/pkg/util"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"sync"

	"github.com/rameshsunkara/go-rest-api-example/internal/config"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
)

var runOnce sync.Once

func Init(serviceInfo *models.ServiceInfo, manager db.DataManager) {
	config := config.GetConfig()
	port := config.GetString("server.port")
	runOnce.Do(func() {
		r := WebRouter(serviceInfo, manager)
		r.Run(":" + port)
	})
}

func WebRouter(svcInfo *models.ServiceInfo, dbMgr db.DataManager) (router *gin.Engine) {
	ginMode := gin.ReleaseMode
	if util.IsDevMode(svcInfo.Environment) {
		ginMode = gin.DebugMode
		gin.ForceConsoleColor()
	}
	gin.SetMode(ginMode)

	// Middleware
	router = gin.Default()
	pprof.Register(router) // TODO: Add debug routes only for Admins /debug/*
	// TODO: Enforce there is authorization information with applicable requests
	// TODO: log everything from gin in json

	// Routes

	// Routes - Status Check
	status := controllers.NewStatusController(svcInfo, dbMgr)
	router.GET("/status", status.CheckStatus) // /status

	// Dependencies for controllers
	d, _ := dbMgr.Database()
	orders := db.NewOrderDataService(d)

	// Routes - Seed DB
	if util.IsDevMode(svcInfo.Environment) {
		seed := controllers.NewSeedController(orders)
		router.POST("/seedDB", seed.SeedDB) // /seedDB
	}

	// Routes - Orders
	v1 := router.Group("/api/v1")
	{
		ordersGroup := v1.Group("orders")
		{
			orders := controllers.NewOrdersController(orders)
			ordersGroup.GET("", orders.GetAll)            // api/v1/orders
			ordersGroup.GET("/:id", orders.GetById)       // api/v1/orders/:id
			ordersGroup.POST("", orders.Post)             // api/v1/orders
			ordersGroup.PUT("", orders.Post)              // api/v1/orders
			ordersGroup.DELETE("/:id", orders.DeleteById) // api/v1/orders/:id
		}
	}

	// Routes - Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return
}
