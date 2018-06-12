// Package http implements the HTTP interface.
package http

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/sirupsen/logrus"
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

func executeTemplate(w http.ResponseWriter, r *http.Request, layout string, layoutData interface{}, name string, data interface{}) {
	logger := middleware.Logger(r)

	// Render template
	var b bytes.Buffer
	if err := templates.ExecuteTemplate(&b, name, data); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to execute layout")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Render layout with template as content
	if err := templates.ExecuteTemplate(w, layout, struct {
		Content template.HTML
		Data    interface{}
	}{
		Content: template.HTML(b.String()),
		Data:    layoutData,
	}); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to execute layout")

		w.WriteHeader(http.StatusInternalServerError)
	}
}
