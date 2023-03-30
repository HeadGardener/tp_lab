package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/HeadHardener/tp_lab/configs"
	"github.com/HeadHardener/tp_lab/internal/app/handlers"
	"github.com/HeadHardener/tp_lab/internal/app/repositories"
	"github.com/HeadHardener/tp_lab/internal/app/services"
	"github.com/HeadHardener/tp_lab/internal/pkg/server"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	"time"
)

var confPath = flag.String("conf-path", "./configs/.env", "path to config env")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(fmt.Sprintf("error whilr creating logger: %s", err.Error()))
	}

	dbconfig, err := configs.NewDBConfig(*confPath)
	if err != nil {
		logger.Fatal(fmt.Sprintf("unable to read config file, error: %s", err.Error()))
	}

	db, err := repositories.NewDB(*dbconfig)
	if err != nil {
		logger.Fatal(fmt.Sprintf("unable to make up conn with db, error: %s", err.Error()))
	}

	repository := repositories.NewRepository(db)
	service := services.NewService(repository)
	handler := handlers.NewHandler(service)

	srvconfig, err := configs.NewServerConfig(*confPath)
	if err != nil {
		logger.Fatal(fmt.Sprintf("unable to read config file, error: %s", err.Error()))
	}

	srv := &server.Server{}

	go func() {
		if err := srv.Run(srvconfig.ServerPort, handler.InitRoutes()); err != nil {
			logger.Error(fmt.Sprintf("error occurring while running server, err:%s", err.Error()))
		}
	}()

	logger.Info("server start working")
	<-ctx.Done()

	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("server forced to shutdown: %s", err.Error()))
	}

	if err := db.Close(); err != nil {
		logger.Error(fmt.Sprintf("db connection forced to shutdown: %s", err.Error()))
	}

	logger.Info("server exiting")
}
