package main

import (
	"context"
	"fmt"
	"main/internal"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	zapLogger "go.uber.org/zap"
)

const maxTimeoutTime = 5

type Server struct {
	port    string
	handler *internal.Handler
	logger  *zapLogger.Logger
}

// NewServer creates a new instance of Server.
func NewServer(port string, handler *internal.Handler, logger *zapLogger.Logger) *Server {
	return &Server{
		port:    port,
		handler: handler,
		logger:  logger,
	}
}

// Run starts the HTTP server and handles graceful shutdown.
func (s *Server) Run() {
	mux := http.NewServeMux()
	s.handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              s.port,
		Handler:           mux,
		ReadHeaderTimeout: maxTimeoutTime * time.Second,
	}

	go func() {
		s.logger.Info(fmt.Sprintf("Starting server on %s", s.port))

		// Start the server and log any errors.
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal("Server failed to start", zapLogger.Error(err))
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal.
	<-stopChan
	s.logger.Info("Shutting down server...")

	// Create a context with a timeout for graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeoutTime*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server.
	if err := server.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown", zapLogger.Error(err))
	}

	s.logger.Info("Server exited properly")
}
