// Package http implements the HTTP interface.
package http

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"github.com/gobuffalo/packr"
)

// key is used within contexts as a key.
type key string

// Keys used by the middleware defined in this package.
const (
	keyVisitor key = "visitor"
)

// Assets and templates loaded using packr.
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

func executeTemplate(w io.Writer, layout string, layoutData interface{}, name string, data interface{}) error {
	// Render template
	var b bytes.Buffer
	if err := templates.ExecuteTemplate(&b, name, data); err != nil {
		return err
	}

	// Render layout with template as content
	return templates.ExecuteTemplate(w, layout, struct {
		Content template.HTML
		Data    interface{}
	}{
		Content: template.HTML(b.String()),
		Data:    layoutData,
	})
}
