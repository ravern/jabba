// Command jabba runs the web server.
package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/ravernkoh/jabba/bolt"
	"github.com/ravernkoh/jabba/http"
	"github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	godotenv.Load()
}

func main() {
	// Get environment variables
	var (
		development                   = os.Getenv("DEVELOPMENT") != ""
		hostname                      = os.Getenv("HOSTNAME")
		port                          = os.Getenv("PORT")
		authSecret                    = os.Getenv("AUTH_SECRET")
		cookieHashKey                 = os.Getenv("COOKIE_HASH_KEY")
		cookieBlockKey                = os.Getenv("COOKIE_BLOCK_KEY")
		googleClientID                = os.Getenv("GOOGLE_CLIENT_ID")
		googleClientSecret            = os.Getenv("GOOGLE_CLIENT_SECRET")
		databasePath                  = os.Getenv("DATABASE_PATH")
		databaseVisitCountIntervalStr = os.Getenv("DATABASE_VISIT_COUNT_INTERVAL")
	)
	if authSecret == "" {
		panic("AUTH_SECRET must be specified")
	}
	if cookieHashKey == "" {
		panic("COOKIE_HASH_KEY must be specified")
	}
	if cookieBlockKey == "" {
		panic("COOKIE_BLOCK_KEY must be specified")
	}
	if googleClientID == "" {
		panic("GOOGLE_CLIENT_ID must be specified")
	}
	if googleClientSecret == "" {
		panic("GOOGLE_CLIENT_SECRET must be specified")
	}
	var databaseVisitCountInterval int
	if databaseVisitCountIntervalStr == "" {
		databaseVisitCountInterval = 10
	} else {
		var err error
		databaseVisitCountInterval, err = strconv.Atoi(databaseVisitCountIntervalStr)
		if err != nil {
			panic("DATABASE_VISIT_COUNT_INTERVAL must be an integer")
		}
	}

	// Create the logger
	logger := logrus.New()
	if !development {
		logger.Formatter = &logrus.JSONFormatter{}
	}

	// Open the database connection
	database := bolt.Database{
		Path:               databasePath,
		VisitCountInterval: time.Duration(databaseVisitCountInterval) * time.Second,
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
		Port:               ":" + port,
		Hostname:           hostname,
		AuthSecret:         authSecret,
		CookieHashKey:      cookieHashKey,
		CookieBlockKey:     cookieBlockKey,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		Logger:             logger,
		Database:           &database,
	}
	// This is a lie!
	logger.Infof("server started listening on port %s", port)
	err := server.Listen()
	logger.WithFields(logrus.Fields{
		"err": err,
	}).Error("server quit")

	// Always exit with an error code since a web server should technically
	// run forever until an error occurs.
	os.Exit(1)
}
