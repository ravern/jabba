package http

import "net/http"

// Landing renders the landing page.
func (s *Server) Landing(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "landing.html", nil)
}

// notFound renders the 404 page.
func notFound(w http.ResponseWriter) {
	templates.ExecuteTemplate(w, "404.html", nil)
}

// internalServerError renders the 500 page.
func internalServerError(w http.ResponseWriter) {
	templates.ExecuteTemplate(w, "500.html", nil)
}
