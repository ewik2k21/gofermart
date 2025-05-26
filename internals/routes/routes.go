package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gofermart/cmd/server"
	"gofermart/internals/handlers"
	"gofermart/internals/interfaces"
	"gofermart/middleware"
)

func RegisterUserRoutes(server server.IGinServer, userHandler *handlers.UserHandler) {
	server.RegisterGroupRoute("api/v1/", []interfaces.RouteDefinition{
		{Method: "POST", Path: "/user/register", Handler: userHandler.Register},
	}, func(ctx *gin.Context) {
		logrus.Info("Request on %s", ctx.Request.URL.Path)
	})

	server.RegisterGroupRoute("api/v1", []interfaces.RouteDefinition{
		{Method: "GET", Path: "/user/get_id", Handler: userHandler.GetId},
	}, middleware.AuthMiddleware())

}
