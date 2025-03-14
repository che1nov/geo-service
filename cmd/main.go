package main

import (
	app2 "geo-service/cmd/app"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	app := app2.NewApp(logger)
	app.Run()
}
