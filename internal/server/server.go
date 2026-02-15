package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/internal/config"
	"github.com/bogdanutanu/go-rest-api-example/internal/db"
	"github.com/bogdanutanu/go-rest-api-example/internal/handlers"
	"github.com/bogdanutanu/go-rest-api-example/internal/middleware"
	"github.com/bogdanutanu/go-rest-api-example/internal/utilities"
	"github.com/bogdanutanu/go-rest-api-example/pkg/flightrecorder"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/bogdanutanu/go-rest-api-example/pkg/mongodb"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	shutdownTimeoutSeconds   = 5
	readHeaderTimeoutSeconds = 60
)

type Server struct {
	Router *gin.Engine
}

// Start manages the HTTP server lifecycle with graceful shutdown
// This function blocks until the server shuts down or an error occurs.
func Start(ctx context.Context, svcEnv *config.ServiceEnvConfig, lgr logger.Logger, dbMgr mongodb.MongoManager) error {
	router, err := WebRouter(svcEnv, lgr, dbMgr)
	if err != nil {
		return err
	}

	// Log registered routes
	lgr.Info().Msg("Registered routes")
	for _, item := range router.Routes() {
		lgr.Info().Str("method", item.Method).Str("path", item.Path).Send()
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:              ":" + svcEnv.Port,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeoutSeconds * time.Second,
	}

	// Channel to capture server startup errors
	serverErrors := make(chan error, 1)

	// Start server in a single goroutine (managed by this function)
	go func() {
		lgr.Info().Str("port", svcEnv.Port).Msg("Starting server")
		serverErrors <- srv.ListenAndServe()
	}()

	// Block and wait for either shutdown signal or server error
	select {
	case serverErr := <-serverErrors:
		if !errors.Is(serverErr, http.ErrServerClosed) {
			return fmt.Errorf("server failed: %w", serverErr)
		}
		return nil
	case <-ctx.Done():
		lgr.Info().Msg("Shutdown signal received, stopping server...")

		// Graceful shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)
		defer cancel()

		if shutdownErr := srv.Shutdown(shutdownCtx); shutdownErr != nil {
			lgr.Error().Err(shutdownErr).Msg("Server forced to shutdown")
			return shutdownErr
		}

		lgr.Info().Msg("Server shutdown gracefully")
		return nil
	}
}

func WebRouter(svcEnv *config.ServiceEnvConfig, lgr logger.Logger, dbMgr mongodb.MongoManager) (*gin.Engine, error) {
	ginMode := gin.ReleaseMode
	if utilities.IsDevMode(svcEnv.Environment) {
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

	// Initialize flight recorder for slow request tracing (if enabled)
	var fr *flightrecorder.Recorder
	if svcEnv.EnableTracing {
		fr = flightrecorder.NewDefault(lgr)
	}
	router.Use(middleware.RequestLogMiddleware(lgr, fr))

	internalAPIGrp := router.Group("/internal")
	internalAPIGrp.Use(middleware.InternalAuthMiddleware()) // use special auth middleware to handle internal employees
	pprof.RouteRegister(internalAPIGrp, "pprof")
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	status, sHandlerErr := handlers.NewStatusHandler(lgr, dbMgr)
	if sHandlerErr != nil {
		return nil, sHandlerErr
	}
	router.GET("/healthz", status.CheckStatus)

	d := dbMgr.Database()
	ordersRepo, ordersRepoErr := db.NewOrdersRepo(lgr, d)
	if ordersRepoErr != nil {
		return nil, ordersRepoErr
	}

	// This is a dev mode only endpoint (route) to seed the local db
	if utilities.IsDevMode(svcEnv.Environment) {
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
	return router, nil
}
