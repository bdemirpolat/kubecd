package main

import (
	"os"

	"github.com/bdemirpolat/kubecd/pkg/application"
	"github.com/bdemirpolat/kubecd/pkg/application/k8apply"
	"github.com/bdemirpolat/kubecd/pkg/database"
	"github.com/bdemirpolat/kubecd/pkg/logger"
	"github.com/urfave/cli/v2"
)

func main() {
	// init logger
	l, err := logger.Init()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = l.Sync()
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
	applicationCLIHandler := application.NewCLIHandler(applicationService)
	app := &cli.App{}
	applicationCLIHandler.RegisterCommands(app)

	err = app.Run(os.Args)
	if err != nil {
		logger.SugarLogger.Fatal(err)
	}
}
