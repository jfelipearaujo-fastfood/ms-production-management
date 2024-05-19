package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/jfelipearaujo-org/ms-production-management/internal/adapter/cloud"
	"github.com/jfelipearaujo-org/ms-production-management/internal/adapter/database"
	"github.com/jfelipearaujo-org/ms-production-management/internal/environment"
	"github.com/jfelipearaujo-org/ms-production-management/internal/handler/get_by_id"
	"github.com/jfelipearaujo-org/ms-production-management/internal/handler/get_by_state"
	"github.com/jfelipearaujo-org/ms-production-management/internal/handler/health"
	"github.com/jfelipearaujo-org/ms-production-management/internal/handler/update"
	"github.com/jfelipearaujo-org/ms-production-management/internal/provider/time_provider"
	"github.com/jfelipearaujo-org/ms-production-management/internal/repository/order_production"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/create"
	get_by_id_service "github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/get_by_id"
	get_by_state_service "github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/get_by_state"
	update_service "github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/update"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Config                  *environment.Config
	DatabaseService         database.DatabaseService
	QueueService            cloud.QueueService
	UpdateOrderTopicService cloud.TopicService

	Dependency Dependency
}

func NewServer(config *environment.Config) *Server {
	ctx := context.Background()

	cloudConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	if config.CloudConfig.IsBaseEndpointSet() {
		cloudConfig.BaseEndpoint = aws.String(config.CloudConfig.BaseEndpoint)
	}

	databaseService := database.NewDatabase(config)

	timeProvider := time_provider.NewTimeProvider(time.Now)
	orderProductionRepository := order_production.NewOrderProductionRepository(databaseService.GetInstance())

	createOrderProductionService := create.NewService(orderProductionRepository, timeProvider)

	updateOrderTopicService := cloud.NewUpdateOrderTopicService(config.CloudConfig.UpdateOrderTopic, cloudConfig)

	return &Server{
		Config:          config,
		DatabaseService: databaseService,
		QueueService: cloud.NewQueueService(
			config.CloudConfig.OrderProductionQueue,
			cloudConfig,
			createOrderProductionService,
			updateOrderTopicService,
		),
		UpdateOrderTopicService: updateOrderTopicService,
		Dependency: Dependency{
			TimeProvider: timeProvider,

			OrderProductionRepository: orderProductionRepository,

			GetOrderProductionById:    get_by_id_service.NewService(orderProductionRepository),
			GetOrderProductionByState: get_by_state_service.NewService(orderProductionRepository),
			UpdateOrderProduction:     update_service.NewService(orderProductionRepository, timeProvider),

			UpdateOrderTopicService: updateOrderTopicService,
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

	group := e.Group(fmt.Sprintf("/api/%s", s.Config.ApiConfig.ApiVersion))

	s.registerOrderProductionHandlers(group)

	return e
}

func (server *Server) registerHealthCheck(e *echo.Echo) {
	healthHandler := health.NewHandler(server.DatabaseService)

	e.GET("/health", healthHandler.Handle)
}

func (s *Server) registerOrderProductionHandlers(e *echo.Group) {
	getOrderProductionByIdHandler := get_by_id.NewHandler(s.Dependency.GetOrderProductionById)
	getOrderProductionByStateHandler := get_by_state.NewHandler(s.Dependency.GetOrderProductionByState)
	updateOrderProductionHandler := update.NewHandler(s.Dependency.UpdateOrderProduction, s.Dependency.UpdateOrderTopicService)

	e.GET("/production/:id", getOrderProductionByIdHandler.Handle)
	e.GET("/production", getOrderProductionByStateHandler.Handle)
	e.PATCH("/production/:id", updateOrderProductionHandler.Handle)
}
