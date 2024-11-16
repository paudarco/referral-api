package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/paudarco/referral-api/internal/config"
	"github.com/paudarco/referral-api/internal/handler"
	"github.com/paudarco/referral-api/internal/repository"
	"github.com/paudarco/referral-api/internal/server"
	"github.com/paudarco/referral-api/internal/service"
	"github.com/paudarco/referral-api/internal/storage"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	pool, err := repository.NewPostresPool(cfg.DB)
	if err != nil {
		logrus.Fatalf("error creating pool: %s", err.Error())
	}

	storage := storage.NewStorage()
	repos := repository.NewRepository(pool)
	services := service.NewService(repos, *cfg, storage)
	handler := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.Server, handler.InitRouters()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	fmt.Println("Referrap api started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Println("Stopping Referral api...")

	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("error while shutting down: %s\n", err.Error())
	}

	pool.Close()

}
