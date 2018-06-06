package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/ravernkoh/jabba/http"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	// Get environment variables
	var (
		development = os.Getenv("DEVELOPMENT") != ""
		port        = os.Getenv("PORT")
	)

	// Create the logger
	logger := logrus.New()
	if !development {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	logger.WithFields(logrus.Fields{
		"development": development,
	}).Info("main: created logger")

	// Start up the server
	server := http.Server{
		Port:   port,
		Logger: logger,
	}
	logger.WithFields(logrus.Fields{
		"port": port,
	}).Info("main: created server")
	if err := server.Listen(); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("main: server quit")
		os.Exit(1)
	}
}
