package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jfelipearaujo-org/ms-production-management/internal/adapter/database"
	"github.com/jfelipearaujo-org/ms-production-management/internal/environment"
	"github.com/jfelipearaujo-org/ms-production-management/internal/handler/health"
	"github.com/jfelipearaujo-org/ms-production-management/internal/provider/time_provider"
	"github.com/jfelipearaujo-org/ms-production-management/internal/repository/order_production"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Config          *environment.Config
	DatabaseService database.DatabaseService

	Dependency Dependency
}

func NewServer(config *environment.Config) *Server {
	databaseService := database.NewDatabase(config)

	timeProvider := time_provider.NewTimeProvider(time.Now)
	orderProductionRepository := order_production.NewOrderProductionRepository(databaseService.GetInstance())

	return &Server{
		Config:          config,
		DatabaseService: databaseService,
		Dependency: Dependency{
			TimeProvider:              timeProvider,
			OrderProductionRepository: orderProductionRepository,
		},
	}
}

func (s *Server) GetHttpServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.ApiConfig.Port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Recover())

	s.registerHealthCheck(e)

	return e
}

func (server *Server) registerHealthCheck(e *echo.Echo) {
	healthHandler := health.NewHandler(server.DatabaseService)

	e.GET("/health", healthHandler.Handle)
}
