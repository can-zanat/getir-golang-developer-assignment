package main

import (
	"fmt"
	"main/config"
	"main/internal"
	"main/store"
	"os"

	logger "github.com/can-zanat/gologger"
	_ "github.com/go-sql-driver/mysql"
)

const serverPort = ":80"

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	// Initialize the application configuration.
	appConfig, err := config.New()
	if err != nil {
		return err
	}

	// Initialize logger with info level.
	loggerInfoLevel := logger.NewWithLogLevel("info")
	defer func() {
		if err := loggerInfoLevel.Sync(); err != nil {
			fmt.Println(err)
		}
	}()

	repository, _ := store.NewStore(appConfig.Mongo)
	service := internal.NewService(repository)
	handler := internal.NewHandler(service)

	// Create and start the HTTP server.
	server := NewServer(serverPort, handler, loggerInfoLevel)
	server.Run()

	return nil
}
