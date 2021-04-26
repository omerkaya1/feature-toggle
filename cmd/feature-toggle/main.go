package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/omerkaya1/feature-toggle/internal"
	"github.com/omerkaya1/feature-toggle/internal/rest"
)

func main() {
	var err error
	// Initialise logger
	log := internal.NewBaseLogger()
	// Init the HTTP server
	server := rest.NewServer(log)

	// Listen for the interrupt signal that will trigger the server shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Create context to propagate global cancellation
	ctx, cancel := context.WithCancel(context.Background())
	go func(cancelFunc context.CancelFunc) {
		<-quit
		cancelFunc()
	}(cancel)

	// Shutting down sequence: the server, the DB connection -> report the exit code
	var exitCode int
	if err = server.Run(ctx, ":8080"); err != nil {
		log.Errorf("failure to gracefully stop the HTTP server: %s", err)
		exitCode = 1
	}
	os.Exit(exitCode)
}
