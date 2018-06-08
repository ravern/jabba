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
	}, "go", "build", "-o", "releases/jabba", "./cmd/jabba"); err != nil {
		return err
	}

	// Clean up static files
	if err := sh.RunV("packr", "clean"); err != nil {
		return err
	}

	return nil
}

const SSH = "ravernkoh@jabba.xyz"

// Deploy copies the production binary onto the server
func Deploy() error {
	// Stop the running service
	if err := sh.RunV("ssh", SSH, "sudo systemctl stop jabba"); err != nil {
		return err
	}

	// Copy it onto the server
	if err := sh.RunV("scp", "releases/jabba", SSH+":~/jabba/jabba"); err != nil {
		return err
	}

	// Start the service
	if err := sh.RunV("ssh", SSH, "sudo systemctl start jabba"); err != nil {
		return err
	}

	return nil
}
