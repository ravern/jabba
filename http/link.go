package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"

	"github.com/ravernkoh/jabba/errors"
	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/ravernkoh/jabba/model"
)

// Index renders the index page.
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	user := s.User(r)

	links, err := s.Database.GetLinks(user.LinkSlugs)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": user.Username,
			"err":      err,
		}).Error("failed to fetch links of user")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.WithFields(logrus.Fields{
		"username": user.Username,
	}).Info("fetched links of user")

	flash, _ := RetrieveFlash(w, r)

	executeTemplate(w, r, "layout.html", nil, "index.html", struct {
		Flash    Flash
		Hostname string
		Links    []*model.Link
	}{
		Flash:    flash,
		Hostname: s.Hostname,
		Links:    links,
	})
}

// RedirectSlug redirects to the corresponding page from the slug.
func (s *Server) RedirectSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	link, err := s.Database.GetLink(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	s.Database.IncrementLinkCount(link)

	http.Redirect(w, r, link.URL, http.StatusFound)
}

// ShortenURL shortens the URL and creates the resulting link.
func (s *Server) ShortenURL(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	user := s.User(r)

	url := r.FormValue("url")
	link, err := model.NewLink(url)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to create link")

		f := Flash{Failure: "Invalid URL given."}
		f.Save(w)

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err := s.Database.CreateLink(link, user); err != nil {
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
func (s *Server) DeleteLink(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	user := s.User(r)

	slug := chi.URLParam(r, "slug")

	if err := s.Database.DeleteLink(slug, user); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to delete link")

		switch err.(errors.Error).Type {
		case errors.Unauthorized:
			f := Flash{Failure: "You can't delete the link!"}
			f.Save(w)
			http.Redirect(w, r, "/", http.StatusFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
