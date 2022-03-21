package main

import (
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
	defer l.Sync()

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

	err = app.Listen(":3001")
	if err != nil {
		logger.SugarLogger.Fatal(err)
	}

}
