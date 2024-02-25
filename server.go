package main

import (
	"context"
	"fmt"
	zapLogger "go.uber.org/zap"
	"main/internal"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//TODO: vakit kalırsa zap logger değiştir. 43 ve 57 de sorun oluyordu vakit kalırsa mutlaka değiştir.

type Server struct {
	port    string
	handler *internal.Handler
	logger  *zapLogger.Logger
}

func NewServer(port string, handler *internal.Handler, logger *zapLogger.Logger) *Server {
	return &Server{
		port:    port,
		handler: handler,
		logger:  logger,
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()
	s.handler.RegisterRoutes(mux) // Assuming a method to register routes

	server := &http.Server{
		Addr:    s.port,
		Handler: mux,
	}

	go func() {
		s.logger.Info(fmt.Sprintf("Starting server on %s", s.port))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal("Server failed to start", zapLogger.Error(err))
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stopChan
	s.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown", zapLogger.Error(err))
	}

	s.logger.Info("Server exited properly")
}
