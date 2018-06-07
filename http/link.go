package http

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ravernkoh/jabba/model"
)

// Index renders the index page.
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	user, err := s.Database.FetchUser("johnsmith")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	links, err := s.Database.FetchLinks(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	executeTemplate(w, "layout.html", nil, "index.html", struct {
		Links []*model.Link
	}{
		Links: links,
	})
}

// Redirect redirects to the corresponding page from the slug.
func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	link, err := s.Database.FetchLink(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link.URL.String(), http.StatusFound)
}
