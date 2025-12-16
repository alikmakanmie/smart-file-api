package config

import (
	"os"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Set output to file and stdout
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}

	// Set log format
	Log.SetFormatter(&logrus.JSONFormatter{})

	// Set log level
	Log.SetLevel(logrus.InfoLevel)

	Log.Info("Logger initialized successfully")
}
