// Command jabba runs the web server.
package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/ravernkoh/jabba/bolt"
	"github.com/ravernkoh/jabba/http"
	"github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	// Get environment variables
	var (
		development  = os.Getenv("DEVELOPMENT") != ""
		hostname     = os.Getenv("HOSTNAME")
		port         = os.Getenv("PORT")
		databasePath = os.Getenv("DATABASE_PATH")
	)

	// Create the logger
	logger := logrus.New()
	if !development {
		logger.Formatter = &logrus.JSONFormatter{}
	}

	// Open the database connection
	database := bolt.Database{
		Path: databasePath,
	}
	if err := database.Open(); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("failed to open database at %s", databasePath)
		os.Exit(1)
	}
	defer database.Close()
	logger.Infof("opened database at %s", databasePath)

	// Start up the server
	server := http.Server{
		Port:     port,
		Hostname: hostname,
		Logger:   logger,
		Database: &database,
	}
	// This is a lie!
	logger.Infof("server started listening on %s", port)
	err := server.Listen()
	logger.WithFields(logrus.Fields{
		"err": err,
	}).Error("server quit")

	// Always exit with an error code since a web server should technically
	// run forever until an error occurs.
	os.Exit(1)
}
