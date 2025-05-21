package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	server "gofermart/cmd/server"
	"gofermart/config"
	"os"
)

func Execute() {
	ctx := context.Background()
	config.LoadEnviroment()
	ginServer, err := server.NewGinServer(ctx, os.Getenv(config.AppPort))
	if err != nil {
		logrus.Fatalf("Failed to create gin server: %v", err)
		return
	}
	err = ginServer.Start(ctx)
	if err != nil {
		logrus.Fatalf("Failed to start gin server: %v", err)
		return
	}

}
