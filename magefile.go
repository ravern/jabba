//+build mage

package main

import "github.com/magefile/mage/sh"

var Default = Development

var Aliases = map[string]interface{}{
	"dev":  Development,
	"prod": Production,
}

// Development starts the development server
func Development() error {
	return sh.RunV("go", "run", "cmd/jabba/main.go")
}

// Production builds the production binary
func Production() error {
	// Generate static files
	if err := sh.RunV("packr"); err != nil {
		return err
	}

	// Build the Linux binary
	if err := sh.RunWith(map[string]string{
		"GOOS":   "linux",
		"GOARCH": "amd64",
	}, "go", "build", "-o", "releases/jabba"); err != nil {
		return err
	}

	// Clean up static files
	if err := sh.RunV("packr", "clean"); err != nil {
		return err
	}

	return nil
}
