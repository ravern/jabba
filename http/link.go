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

// Redirect redirects to the corresponding page from the slug.
//
// TODO: Add authentication
func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	link, err := s.Database.FetchLink(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}

// Shorten shortens the URL and creates the resulting link.
//
// TODO: Support both Visitor and User
func (s *Server) Shorten(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	visitor := s.Visitor(r)

	url := r.FormValue("url")
	link := model.NewLink(url)

	if err := s.Database.CreateVisitorLink(link, visitor); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to create link")

		// TODO: Improve flash message
		setFlash(w, "Failed to shorten link!")

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	logrus.WithFields(logrus.Fields{
		"slug": link.Slug,
	}).Info("created link")

	http.Redirect(w, r, "/", http.StatusFound)
}
