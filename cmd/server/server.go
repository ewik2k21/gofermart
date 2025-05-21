package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ginServer struct {
	engine   *gin.Engine
	server   *http.Server
	startErr chan error
}

func NewGinServer(ctx context.Context, httpAddress string) (*ginServer, error) {
	engine := gin.Default()
	//engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
		logrus.Info("shutting down server...")
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
		logrus.Errorf("server shutdown faailed : %v", err)
		return err
	}

	logrus.Info("Server stopped")
	return nil
}
