package http

import "net/http"

// notFound renders the 404 page.
func notFound(w http.ResponseWriter) {
	executeTemplate(w, "layout.html", nil, "errors/404.html", nil)
}

// internalServerError renders the 500 page.
func internalServerError(w http.ResponseWriter) {
	executeTemplate(w, "layout.html", nil, "errors/500.html", nil)
}
