// Command jabba runs the web server.
package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/ravernkoh/jabba/bolt"
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
		development  = os.Getenv("DEVELOPMENT") != ""
		port         = os.Getenv("PORT")
		databasePath = os.Getenv("DATABASE_PATH")
	)

	// Create the logger
	logger := logrus.New()
	if !development {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	logger.WithFields(logrus.Fields{
		"development": development,
	}).Info("main: created logger")

	// Open the database connection
	database := bolt.Database{
		Path: databasePath,
	}
	if err := database.Open(); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("main: failed to open database connection")
		os.Exit(1)
	}
	logger.WithFields(logrus.Fields{
		"path": databasePath,
	}).Info("main: opened database connection")

	// Start up the server
	server := http.Server{
		Port:   port,
		Logger: logger,
	}
	logger.WithFields(logrus.Fields{
		"port": port,
	}).Info("main: server started listening")
	if err := server.Listen(); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("main: server quit")
		os.Exit(1)
	}
}
