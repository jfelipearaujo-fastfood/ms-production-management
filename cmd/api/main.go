package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/jfelipearaujo-org/ms-production-management/internal/adapter/cloud"
	"github.com/jfelipearaujo-org/ms-production-management/internal/environment"
	"github.com/jfelipearaujo-org/ms-production-management/internal/environment/loader"
	"github.com/jfelipearaujo-org/ms-production-management/internal/server"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/logger"
)

func init() {
	var err error
	time.Local, err = time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()

	loader := loader.NewLoader()

	var config *environment.Config
	var err error

	if len(os.Args) > 1 && os.Args[1] == "local" {
		slog.InfoContext(ctx, "loading environment from .env file")
		config, err = loader.GetEnvironmentFromFile(ctx, ".env")
	} else {
		config, err = loader.GetEnvironment(ctx)
	}

	if err != nil {
		slog.ErrorContext(ctx, "error loading environment", "error", err)
		panic(err)
	}

	logger.SetupLog(config)

	cloudConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	if config.CloudConfig.IsBaseEndpointSet() {
		cloudConfig.BaseEndpoint = aws.String(config.CloudConfig.BaseEndpoint)
	}

	secret := cloud.NewSecretService(cloudConfig)

	dbUrl, err := secret.GetSecret(ctx, config.DbConfig.UrlSecretName)
	if err != nil {
		slog.ErrorContext(ctx, "error getting secret", "secret_name", config.DbConfig.UrlSecretName, "error", err)
		panic(err)
	}

	config.DbConfig.Url = dbUrl

	server := server.NewServer(config)

	if err := server.UpdateOrderTopicService.UpdateTopicArn(ctx); err != nil {
		slog.ErrorContext(ctx, "error updating update order topic url", "error", err)
		panic(err)
	}

	if err := server.QueueService.UpdateQueueUrl(ctx); err != nil {
		slog.ErrorContext(ctx, "error updating queue url", "error", err)
		panic(err)
	}

	go func(ctx context.Context) {
		for {
			server.QueueService.ConsumeMessages(ctx)
		}
	}(ctx)

	httpServer := server.GetHttpServer()

	go func() {
		slog.InfoContext(ctx, "🚀 Server started", "address", httpServer.Addr)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "http server error", "error", err)
			panic(err)
		}
		slog.InfoContext(ctx, "http server stopped serving requests")
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sc

	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	if err := httpServer.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "error while trying to shutdown the server", "error", err)
	}
	slog.InfoContext(ctx, "graceful shutdown completed ✅")
}
