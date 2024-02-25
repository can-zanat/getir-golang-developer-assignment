package main

import (
	"fmt"
	logger "github.com/can-zanat/gologger"
	_ "github.com/go-sql-driver/mysql"
	"main/config"
	"main/internal"
	"main/store"
	"os"
)

const serverPort = ":80"

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	appConfig, err := config.New()
	if err != nil {
		return err
	}

	loggerInfoLevel := logger.NewWithLogLevel("info")
	defer func() {
		if err := loggerInfoLevel.Sync(); err != nil {
			fmt.Println(err)
		}
	}()

	repository := store.NewStore(appConfig.Mongo)
	service := internal.NewService(repository)
	handler := internal.NewHandler(service)

	server := NewServer(serverPort, handler, loggerInfoLevel)
	server.Run()

	return nil
}
