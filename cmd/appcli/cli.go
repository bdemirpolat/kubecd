package main

import (
	application2 "github.com/bdemirpolat/kubecd/application"
	"github.com/bdemirpolat/kubecd/application/k8apply"
	"github.com/bdemirpolat/kubecd/database"
	"github.com/bdemirpolat/kubecd/logger"
	"os"

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
	applicationRepo := application2.NewRepo(db)
	applicationService := application2.NewService(applicationRepo)
	applicationCLIHandler := application2.NewCLIHandler(applicationService)
	app := &cli.App{Name: "kubecd", Usage: "kubecd command [command options] [arguments...]"}
	applicationCLIHandler.RegisterCommands(app)

	err = app.Run(os.Args)
	if err != nil {
		logger.SugarLogger.Fatal(err)
	}
}
