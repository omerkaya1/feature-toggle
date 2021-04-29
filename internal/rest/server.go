package rest

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/omerkaya1/feature-toggle/internal"
)

type (
	logger interface {
		Fatalf(format string, args ...interface{})
		Debugf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Infof(format string, args ...interface{})
	}
	storage interface {
		GetFeatures(ctx context.Context) ([]internal.Feature, error)
		GetFeatureByName(ctx context.Context, name string) (internal.Feature, error)
		GetUserFeatures(ctx context.Context, user string) ([]internal.Feature, error)
		GetUserFeaturesByStatus(ctx context.Context, customer string, active bool) ([]internal.Feature, error)
		CreateFeature(ctx context.Context, inverted, active bool, displayName, techName, description string, customerIDs []string, expires time.Time) (int64, error)
		UpdateFeature(ctx context.Context, techName, description, displayName string, ed time.Time, status bool) error
		DeleteFeature(ctx context.Context, techName string) error
	}
	Server struct {
		log   logger
		store storage
	}
)

var reqLogCfg = middleware.LoggerConfig{
	Format: "${time_rfc3339_nano} ${user_agent} ${remote_ip} ${host} ${method} ${uri} ${latency_human} ${status} ${error}\n",
}

func NewServer(log logger, store storage) *Server {
	return &Server{
		log:   log,
		store: store,
	}
}

func (s *Server) Run(ctx context.Context, addr string) error {
	e := echo.New()
	// Apply middleware
	e.Use(middleware.Recover(), middleware.LoggerWithConfig(reqLogCfg), middleware.CORS())
	// Serve frontend
	e.Static("", "./static/dist/ft")
	// API
	api := e.Group("/api/v1")
	api.GET("/ping", func(c echo.Context) error { return c.JSON(200, "pong") })
	api.GET("/features", s.getFeatures)
	api.POST("/features", s.createFeature)
	api.PUT("/features/:name", s.updateFeature)
	api.DELETE("/features/:name", s.deleteFeature)
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
