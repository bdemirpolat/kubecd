package main

import (
	application2 "github.com/bdemirpolat/kubecd/application"
	"github.com/bdemirpolat/kubecd/application/k8apply"
	"github.com/bdemirpolat/kubecd/database"
	"github.com/bdemirpolat/kubecd/logger"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	applicationRepo := application2.NewRepo(db)
	applicationService := application2.NewService(applicationRepo)
	applicationHandler := application2.NewHandler(applicationService)

	app := fiber.New()
	applicationHandler.RegisterHandlers(app)

	err = application2.Loop(applicationRepo)
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
