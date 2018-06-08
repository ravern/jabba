package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"

	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/ravernkoh/jabba/model"
)

// Index renders the index page.
//
// TODO: Support both Visitor and User
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {

	logger := middleware.Logger(r)
	visitor := s.Visitor(r)

	links, err := s.Database.FetchLinks(visitor.LinkSlugs)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"token": visitor.Token,
			"err":   err,
		}).Error("failed to fetch links of visitor")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.WithFields(logrus.Fields{
		"token": visitor.Token,
	}).Info("fetched links of visitor")

	flash, _ := flash(w, r)

	if err := executeTemplate(w, "layout.html", nil, "index.html", struct {
		Flash    string
		Hostname string
		Links    []*model.Link
	}{
		Flash:    flash,
		Hostname: s.Hostname,
		Links:    links,
	}); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to execute template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// RedirectSlug redirects to the corresponding page from the slug.
//
// TODO: Add authentication
func (s *Server) RedirectSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	link, err := s.Database.FetchLink(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := s.Database.IncrementLinkCount(link); err != nil {
		middleware.Logger(r).WithFields(logrus.Fields{
			"err": err,
		}).Warn("failed to increment link count")
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}

// ShortenURL shortens the URL and creates the resulting link.
//
// TODO: Support both Visitor and User
func (s *Server) ShortenURL(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	visitor := s.Visitor(r)

	url := r.FormValue("url")
	link, err := model.NewLink(url)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to create link")

		setFlash(w, "Invalid URL given!")

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err := s.Database.CreateVisitorLink(link, visitor); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to create link")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logrus.WithFields(logrus.Fields{
		"slug": link.Slug,
	}).Info("created link")

	http.Redirect(w, r, "/", http.StatusFound)
}

// DeleteLink deletes the link.
//
// TODO: Support both Visitor and User
func (s *Server) DeleteLink(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	visitor := s.Visitor(r)

	slug := chi.URLParam(r, "slug")

	if err := s.Database.DeleteVisitorLink(slug, visitor); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to delete link")

		setFlash(w, "Trying to delete invalid link!")

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
