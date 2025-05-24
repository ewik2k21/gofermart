package provider

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"gofermart/cmd/server"
	"gofermart/internals/handlers"
	"gofermart/internals/repositories"
	"gofermart/internals/routes"
	"gofermart/internals/services"
)

func NewProvider(pool *pgxpool.Pool, server *server.GinServer) {
	// repo db
	// routes server
	userRepo := repositories.NewUserRepository(pool)
	userService := services.NewUserService(*userRepo)
	tokenService := services.NewTokenService()
	userHandler := handlers.NewUserHandler(*userService, *tokenService)
	routes.RegisterUserRoutes(server, userHandler)
}
