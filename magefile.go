//+build mage

package main

import (
	"os"
	"os/exec"
	"strings"

	fswatch "github.com/andreaskoch/go-fswatch"
	"github.com/magefile/mage/sh"
	"github.com/sirupsen/logrus"
)

var Default = Development

var Aliases = map[string]interface{}{
	"dev":  Development,
	"prod": Production,
}

// Development starts the development server
func Development() error {
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Debug("starting development server")

	var (
		first = make(chan struct{}, 1)
		cmd   *exec.Cmd
	)
	first <- struct{}{}

	w := fswatch.NewFolderWatcher(".", true, func(path string) bool {
		return !(strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, ".env") ||
			strings.HasSuffix(path, ".html") ||
			strings.HasSuffix(path, ".mod"))
	}, 1)

	w.Start()
	for w.IsRunning() {
		select {
		case <-first:
		case <-w.ChangeDetails():
		}

		if cmd != nil {
			logrus.Debug("restarting server")
			if err := cmd.Process.Kill(); err != nil {
				return err
			}
		}

		logrus.Debug("building the binary")
		if err := sh.RunV("go", "build", "-o", "releases/jabba", "./cmd/jabba"); err != nil {
			logrus.Error("failed to start server: ", err)
			continue
		}

		logrus.Debug("running the binary")
		cmd = exec.Command("releases/jabba")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			logrus.Error("failed to start server: ", err)
		}
	}

	return nil
}

// Production builds the production binary
func Production() error {
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Debug("generating static files")
	if err := sh.RunV("packr"); err != nil {
		return err
	}

	logrus.Debug("building the binary")
	if err := sh.RunV("go", "build", "-o", "releases/jabba", "./cmd/jabba"); err != nil {
		return err
	}

	logrus.Debug("cleaning static files")
	return sh.RunV("packr", "clean")
}
