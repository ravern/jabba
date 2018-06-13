package http

import (
	"net/http"
)

// notFound renders the 404 page.
func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"error.css",
	}, nil, "error.html", map[string]interface{}{
		"CurrentUsername": s.currentUsername(r),
		"Message":         "404 Not Found",
	})
}

// unauthorized renders the 401 page.
func (s *Server) unauthorized(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"error.css",
	}, nil, "error.html", map[string]interface{}{
		"CurrentUsername": s.currentUsername(r),
		"Message":         "401 Unauthorized",
	})
}

// internalServerError renders the 500 page.
func (s *Server) internalServerError(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"error.css",
	}, nil, "error.html", map[string]interface{}{
		"CurrentUsername": s.currentUsername(r),
		"Message":         "500 Internal Server Error",
	})
}
