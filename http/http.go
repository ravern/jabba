package http

import (
	"html/template"
	"net/http"

	"github.com/gobuffalo/packr"
)

var (
	assets    http.FileSystem
	templates *template.Template
)

func init() {
	// Load assets
	assets = packr.NewBox("./assets")

	// Load templates
	templates = template.New("")
	box := packr.NewBox("./templates")
	box.Walk(func(name string, f packr.File) error {
		template.Must(templates.New(name).Parse(box.String(name)))
		return nil
	})
}
