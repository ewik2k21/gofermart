package main

import (
	"gofermart/cmd"
	_ "gofermart/docs"
)

// @title gofermartAPI
// @version 1.0
// @description API server for gofermart
// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
