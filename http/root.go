package http

import "net/http"

// Root renders the index page.
func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "layout.html", nil, "index.html", nil)
}
