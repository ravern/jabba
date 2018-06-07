package http

import "net/http"

// notFound renders the 404 page.
func notFound(w http.ResponseWriter) {
	executeTemplate(w, "layout.html", nil, "error.html", struct {
		Message string
	}{
		Message: "404 Not Found",
	})
}

// internalServerError renders the 500 page.
func internalServerError(w http.ResponseWriter) {
	executeTemplate(w, "layout.html", nil, "error.html", struct {
		Message string
	}{
		Message: "500 Internal Server Error",
	})
}
