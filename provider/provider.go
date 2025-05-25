package provider

import (
	"database/sql"
	"gofermart/cmd/server"
	"gofermart/internals/handlers"
	"gofermart/internals/repositories"
	"gofermart/internals/routes"
	"gofermart/internals/services"
)

func NewProvider(db *sql.DB, server server.IGinServer) {
	// repo db
	// routes server
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	tokenService := services.NewTokenService()
	userHandler := handlers.NewUserHandler(userService, tokenService)
	routes.RegisterUserRoutes(server, userHandler)
}
