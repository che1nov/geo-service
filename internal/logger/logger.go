package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)

	return log
}
