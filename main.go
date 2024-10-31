package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ideal-tekno-solusi/sso/api"
	"github.com/ideal-tekno-solusi/sso/bootstrap"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()
	cfg := bootstrap.InitContainer()

	api.RegisterApi(r, cfg)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("failed to listen with error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Warn("shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg.StopDb(ctx)

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("shutdown server with error: %v", err)
	}

	select {
	case <-ctx.Done():
		logrus.Warn("timeout of 5 seconds.")
	}

	logrus.Warn("Server exiting")
}
