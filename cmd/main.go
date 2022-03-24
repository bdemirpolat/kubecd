package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bdemirpolat/kubecd/pkg/application"
	"github.com/bdemirpolat/kubecd/pkg/application/k8apply"
	"github.com/bdemirpolat/kubecd/pkg/database"
	"github.com/bdemirpolat/kubecd/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// init logger
	l, err := logger.Init()
	if err != nil {
		panic(err)
	}
	defer func() {
		syncErr := l.Sync()
		if syncErr != nil {
			log.Fatal(err)
		}
	}()

	// init database
	db, err := database.Init()
	if err != nil {
		logger.SugarLogger.Fatal(err)
	}

	err = k8apply.InitKubeClient()
	if err != nil {
		logger.SugarLogger.Fatal(err)
	}

	// init rest server
	applicationRepo := application.NewRepo(db)
	applicationService := application.NewService(applicationRepo)
	applicationHandler := application.NewHandler(applicationService)

	app := fiber.New()
	applicationHandler.RegisterHandlers(app)

	err = application.Loop(applicationRepo)
	if err != nil {
		logger.SugarLogger.Fatal(err)
	}

	go func() {
		listenErr := app.Listen(":3001")
		if listenErr != nil {
			logger.SugarLogger.Fatal(listenErr)
		}
	}()

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-interruptSignal
}
