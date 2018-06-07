package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"

	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/ravernkoh/jabba/model"
)

// Index renders the index page.
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)

	username := "johnsmith"
	user, err := s.Database.FetchUser(username)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": username,
			"err":      err,
		}).Error("failed to fetch user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.WithFields(logrus.Fields{
		"username": username,
	}).Info("fetched user")

	links, err := s.Database.FetchLinks(user)
	if err != nil {
		logger.Error("failed to fetch links")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("fetched links")

	if err := executeTemplate(w, "layout.html", nil, "index.html", struct {
		Hostname string
		Links    []*model.Link
	}{
		Hostname: s.Hostname,
		Links:    links,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Redirect redirects to the corresponding page from the slug.
func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	link, err := s.Database.FetchLink(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}
