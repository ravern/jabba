// Package http implements the HTTP interface.
package http

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/ravern/jabba/http/middleware"
	"github.com/sirupsen/logrus"
)

// key is used within contexts as a key.
type key string

// Keys used by the middleware defined in this package.
const (
	keyUser key = "user"
	keyLink key = "link"
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
	templates = template.New("").Funcs(map[string]interface{}{
		"join": strings.Join,
	})
	box := packr.NewBox("./templates")
	box.Walk(func(name string, f packr.File) error {
		template.Must(templates.New(name).Parse(box.String(name)))
		return nil
	})

	// Load templates functions
}

func executeTemplate(w http.ResponseWriter, r *http.Request, layout string, styles []string, layoutData interface{}, name string, data interface{}) {
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
	if err := templates.ExecuteTemplate(w, layout, map[string]interface{}{
		"Content": template.HTML(b.String()),
		"Styles":  styles,
		"Data":    layoutData,
	}); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to execute layout")

		w.WriteHeader(http.StatusInternalServerError)
	}
}
