package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gofermart/internals/interfaces"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GinServer struct {
	engine   *gin.Engine
	server   *http.Server
	startErr chan error
}

func NewGinServer(ctx context.Context, httpAddress string) (*GinServer, error) {
	engine := gin.Default()
	//engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	gs := &GinServer{
		engine:   engine,
		server:   &http.Server{Addr: httpAddress, Handler: engine},
		startErr: make(chan error, 1),
	}

	go func() {
		logrus.Infof("start http server at %s", httpAddress)
		if err := gs.server.ListenAndServe(); err != nil {
			logrus.Errorf("listening on %s faailed : %v", httpAddress, err)
			gs.startErr <- err
		}
	}()

	return gs, nil

}
func (gs *GinServer) Start(ctx context.Context) error {
	if gs.server == nil {
		return errors.New("server is nil")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-gs.startErr:
		return err
	case <-quit:
		logrus.Info("shutting down server...")
		return gs.server.Shutdown(ctx)

	case <-ctx.Done():
		logrus.Info("Server shutting down, context cancelled")

		return gs.server.Shutdown(ctx)
	}

}

func (gs *GinServer) ShutDown(ctx context.Context) error {
	ctxShutDown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := gs.server.Shutdown(ctxShutDown); err != nil {
		logrus.Errorf("server shutdown faailed : %v", err)
		return err
	}

	logrus.Info("Server stopped")
	return nil
}

func (gs *GinServer) RegisterGroupRoute(path string, routes []interfaces.RouteDefinition, middleWare ...gin.HandlerFunc) {
	group := gs.engine.Group(path)
	group.Use(middleWare...)
	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Path, route.Handler)
		case "POST":
			group.POST(route.Path, route.Handler)
		case "PUT":
			group.PUT(route.Path, route.Handler)
		case "DELETE":
			group.DELETE(route.Path, route.Handler)
		case "PATCH":
			group.PATCH(route.Path, route.Handler)

		default:
			logrus.Errorf("Invalid https method")
		}
	}
}
