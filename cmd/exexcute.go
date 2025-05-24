package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"gofermart/cmd/server"
	"gofermart/config"
	"gofermart/provider"
	"os"
)

func Execute() {
	ctx := context.Background()
	config.LoadEnv()
	ginServer, err := server.NewGinServer(ctx, os.Getenv(config.AppPort))
	if err != nil {
		logrus.Fatalf("Failed to create gin server: %v", err)
		return
	}

	db, err := config.SetUpDatabase(ctx)
	if err != nil {
		logrus.Fatalf("Failed to setup database: %v", err)
	}

	err = ginServer.Start(ctx)
	if err != nil {
		logrus.Fatalf("Failed to start gin server: %v", err)
		return
	}

	provider.NewProvider(db, ginServer)

}
