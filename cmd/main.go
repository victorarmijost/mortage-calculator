package main

import (
	"os"
	"varmijo/mortage-calculator/pkg/server"
	"varmijo/mortage-calculator/pkg/usecase"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := setLogger()

	framework := echo.New()

	controller := usecase.New(logger)
	server.Register(framework, controller, logger)

	logger.Fatal(framework.Start(":8080"))
}

func setLogger() *logrus.Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.DebugLevel
	}

	logger.SetLevel(level)

	return logger
}
