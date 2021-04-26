package rest

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	logger interface {
		Fatalf(format string, args ...interface{})
		Debugf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Infof(format string, args ...interface{})
	}
	Server struct {
		log logger
	}
)

var reqLogCfg = middleware.LoggerConfig{
	Format: "${time_rfc3339_nano} ${user_agent} ${remote_id} ${host} ${method} ${uri} ${latency_human} ${status} ${error}\n",
}

func NewServer(log logger) *Server {
	return &Server{log: log}
}

func (s *Server) Run(ctx context.Context, addr string) error {
	e := echo.New()
	// Apply middleware
	e.Use(middleware.Recover(), middleware.LoggerWithConfig(reqLogCfg))
	// Ping
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, "pong")
	})
	e.Static("", "./static/dist/ft")
	go func() {
		if err := e.Start(addr); err != nil {
			s.log.Fatalf("the HTTP server stopped: %s", err)
		}
	}()
	// Wait for the context cancellation
	<-ctx.Done()
	// Crete context to gracefully shut down the server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return e.Shutdown(shutdownCtx)
}
