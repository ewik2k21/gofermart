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

type IGinServer interface {
	Start(ctx context.Context) error
	ShutDown(ctx context.Context) error
	RegisterGroupRoute(path string, routes []interfaces.RouteDefinition, middleWare ...gin.HandlerFunc)
	RegisterRoute(method, path string, handler gin.HandlerFunc)
}

type ginServer struct {
	engine   *gin.Engine
	server   *http.Server
	startErr chan error
}

func NewGinServer(ctx context.Context, httpAddress string) (IGinServer, error) {
	engine := gin.Default()

	gs := &ginServer{
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

func (gs *ginServer) Start(ctx context.Context) error {
	if gs.server == nil {
		return errors.New("server is nil")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-gs.startErr:
		return err
	case <-quit:
		logrus.Info("Shutting down server...")
		return gs.server.Shutdown(ctx)

	case <-ctx.Done():
		logrus.Info("Server shutting down, context cancelled")

		return gs.server.Shutdown(ctx)
	}

}

func (gs *ginServer) ShutDown(ctx context.Context) error {
	ctxShutDown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := gs.server.Shutdown(ctxShutDown); err != nil {
		logrus.Errorf("Server shutdown faailed : %v", err)
		return err
	}

	logrus.Info("Server stopped")
	return nil
}
func (gs *ginServer) RegisterRoute(method, path string, handler gin.HandlerFunc) {
	switch method {
	case "GET":
		gs.engine.GET(path, handler)
	case "POST":
		gs.engine.POST(path, handler)
	case "PUT":
		gs.engine.PUT(path, handler)
	case "DELETE":
		gs.engine.DELETE(path, handler)
	case "PATCH":
		gs.engine.PATCH(path, handler)

	default:
		logrus.Errorf("Invalid https method")
	}

}
func (gs *ginServer) RegisterGroupRoute(path string, routes []interfaces.RouteDefinition, middleWare ...gin.HandlerFunc) {
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
