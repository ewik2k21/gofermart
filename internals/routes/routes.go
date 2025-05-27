package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gofermart/cmd/server"
	"gofermart/internals/handlers"
	"gofermart/internals/interfaces"
)

func RegisterUserRoutes(server server.IGinServer, userHandler *handlers.UserHandler) {
	server.RegisterGroupRoute("api/v1/", []interfaces.RouteDefinition{
		{Method: "POST", Path: "/user/register", Handler: userHandler.Register},
		{Method: "POST", Path: "/user/login", Handler: userHandler.Login},
	}, func(ctx *gin.Context) {
		logrus.Infof("Request on %s", ctx.Request.URL.Path)
	})

	//server.RegisterGroupRoute("api/v1", []interfaces.RouteDefinition{
	//	{Method: "GET", Path: "/user/get_id", Handler: userHandler.GetId},
	//}, middleware.AuthMiddleware())

}
