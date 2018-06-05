package main

import (
	"fmt"
	"os"

	"github.com/gobuffalo/packr"
	"github.com/joho/godotenv"
	"github.com/ravernkoh/jabba/http"
	"github.com/sirupsen/logrus"
)

func main() {
	must(godotenv.Load(), 1)

	// Load environment variables
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
	}).Info("created logger")

	// Load static assets
	assets := packr.NewBox("../../assets")
	logger.Info("loaded static assets")

	// Start up the server
	server := http.Server{
		Port:   port,
		Assets: assets,
		Logger: logger,
	}
	logger.WithFields(logrus.Fields{
		"port": port,
	}).Info("server listening...")
	must(server.Listen(), 1)
}

// must will exit the application with the given exit code if the given error
// is non-nil.
func must(err error, code int) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(code)
	}
}
