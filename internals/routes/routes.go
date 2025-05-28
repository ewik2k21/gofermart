package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gofermart/cmd/server"
	"gofermart/internals/handlers"
	"gofermart/internals/interfaces"
	"gofermart/middleware"
)

func RegisterUserRoutes(server server.IGinServer, userHandler *handlers.UserHandler) {
	server.RegisterGroupRoute("api/v1/", []interfaces.RouteDefinition{
		{Method: "GET", Path: "/swagger/*any", Handler: ginSwagger.WrapHandler(swaggerFiles.Handler)},
		{Method: "POST", Path: "/user/register", Handler: userHandler.Register},
		{Method: "POST", Path: "/user/login", Handler: userHandler.Login},
	}, func(ctx *gin.Context) {
		logrus.Infof("Request on %s", ctx.Request.URL.Path)
	})

	server.RegisterGroupRoute("api/v1", []interfaces.RouteDefinition{
		{Method: "POST", Path: "/user/orders", Handler: userHandler.AddOrder},
		{Method: "GET", Path: "/user/orders", Handler: userHandler.GetAllOrders},
		{Method: "GET", Path: "/user/balance", Handler: userHandler.GetBalance},
	}, func(ctx *gin.Context) {
		logrus.Infof("Request on %s", ctx.Request.URL.Path)
	}, middleware.AuthMiddleware())

}
