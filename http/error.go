package http

import (
	"net/http"
)

// notFound renders the 404 page.
func notFound(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"error.css",
	}, nil, "error.html", map[string]interface{}{
		"Message": "404 Not Found",
	})
}

// internalServerError renders the 500 page.
func internalServerError(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"error.css",
	}, nil, "error.html", map[string]interface{}{
		"Message": "500 Internal Server Error",
	})
}
