package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"

	"github.com/ravern/jabba/errors"
	"github.com/ravern/jabba/http/middleware"
	"github.com/ravern/jabba/model"
)

// SetLink sets the link in the context.
//
// Must be placed after SetLogger.
func (s *Server) SetLink(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.Logger(r)

		slug := chi.URLParam(r, "slug")

		l, err := s.Database.GetLink(slug)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"slug": slug,
			}).Warnf("failed to get link")

			w.WriteHeader(http.StatusNotFound)
			return
		}
		logger.WithFields(logrus.Fields{
			"slug": slug,
		}).Info("got link")

		ctx := context.WithValue(r.Context(), keyLink, l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireLink checks if the user is authorized to access the link.
//
// Must be placed after SetLogger, SetUser and SetLink.
func (s *Server) RequireLink(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.Logger(r)
		user := s.User(r)
		link := s.Link(r)

		if _, ok := user.FindLinkSlug(link.Slug); !ok {
			logger.WithFields(logrus.Fields{
				"slug": link.Slug,
			}).Warn("get link unauthorized")

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Link returns the link for the given request.
func (s *Server) Link(r *http.Request) *model.Link {
	return r.Context().Value(keyLink).(*model.Link)
}

// Index renders the index page.
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	user := s.User(r)

	links, err := s.Database.GetLinks(user.LinkSlugs)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": user.Username,
			"err":      err,
		}).Error("failed to get links of user")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.WithFields(logrus.Fields{
		"username": user.Username,
	}).Info("got links of user")

	flash, _ := s.Flash(w, r)

	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"index.css",
	}, nil, "index.html", map[string]interface{}{
		"CurrentUsername": s.currentUsername(r),
		"Flash":           flash,
		"Hostname":        s.Hostname,
		"Links":           links,
	})
}

// Redirect redirects to the corresponding page from the slug.
func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	link := s.Link(r)

	auths, err := s.Database.GetAuths(link.AuthIDs)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to get auths")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	model.SortAuths(auths)

	var authed bool
	for _, auth := range auths {
		switch auth.Method {
		case model.MethodJabba:
			user := s.User(r)
			for _, username := range auth.Values {
				if user.Username == username {
					authed = true
					break
				}
			}
		case model.MethodGoogle:
			a, err := s.Auth(w, r)
			if err != nil {
				logger.Info("redirecting to Google OAuth")

				http.Redirect(w, r, s.googleConfig.AuthCodeURL(link.Slug), http.StatusFound)
				return
			}

			for _, email := range auth.Values {
				if a.Google == email {
					authed = true
					break
				}
			}
		}
	}
	if len(auths) == 0 {
		authed = true
	}

	if !authed {
		logger.Warn("redirect link unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s.Database.IncrementLinkCount(link)
	http.Redirect(w, r, link.URL, http.StatusFound)
}

// CreateLink shortens the URL and creates the resulting link.
func (s *Server) CreateLink(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	user := s.User(r)

	url := r.FormValue("url")
	link, err := model.NewLink(url)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to create link")

		s.SetFlash(w, Flash{Failure: "Could not create link."})
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

	s.SetFlash(w, Flash{Success: "Successfully created link!"})
	http.Redirect(w, r, "/", http.StatusFound)
}

// UpdateLinkForm renders to form to update the link.
func (s *Server) UpdateLinkForm(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	link := s.Link(r)
	flash, _ := s.Flash(w, r)

	auths, err := s.Database.GetAuths(link.AuthIDs)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"slug": link.Slug,
		}).Error("failed to get auths for link")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.executeUpdateLinkFormTemplate(w, r, flash, link, auths)
}

// UpdateLink updates the link.
func (s *Server) UpdateLink(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	user := s.User(r)
	link := s.Link(r)

	slug := link.Slug
	link.Slug = r.FormValue("slug")
	link.Title = r.FormValue("title")
	link.URL = r.FormValue("url")

	// Parse the auths in the form
	if err := r.ParseForm(); err != nil {
		logger.Error("failed to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var f Flash

	var (
		auths       []*model.Auth
		authIDs     = r.PostForm["auth[id]"]
		authMethods = r.PostForm["auth[method]"]
		authValuess = r.PostForm["auth[values]"]
	)
	if len(authMethods) != len(authValuess) {
		logger.Warn("invalid auths given")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for i, method := range authMethods {
		if method == "" {
			continue
		}

		auth, err := model.NewAuth(method, authValuess[i])
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Warn("failed to create auth")

			f.Failure = "Could not update link."

			link.Slug = slug
			s.executeUpdateLinkFormTemplate(w, r, f, link, auths)
			return
		}
		if i < len(authIDs) {
			auth.ID = authIDs[i]
		}

		auths = append(auths, auth)
	}

	if err := link.Validate(); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("failed link validation")

		f.Failure = "Could not update link."

		link.Slug = slug
		s.executeUpdateLinkFormTemplate(w, r, f, link, auths)
		return
	}

	if err := s.Database.UpdateLinkSlug(slug, link, user); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to update link")

		switch err.(errors.Error).Type {
		case errors.AlreadyExists:
			f.Failure = "Slug already exists."
		default:
			f.Failure = "Could not update link."
		}

		link.Slug = slug
		s.executeUpdateLinkFormTemplate(w, r, f, link, auths)
		return
	}

	if err := s.Database.UpdateAuths(auths, link); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to update auths")

		f.Failure = "Could not update link."

		link.Slug = slug
		s.executeUpdateLinkFormTemplate(w, r, f, link, auths)
		return
	}

	s.SetFlash(w, Flash{Success: "Successfully updated link!"})
	http.Redirect(w, r, "/", http.StatusFound)
}

// DeleteLink deletes the link.
func (s *Server) DeleteLink(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	user := s.User(r)
	link := s.Link(r)

	if err := s.Database.DeleteLink(link, user); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("delete link unauthorized")

		switch err.(errors.Error).Type {
		case errors.Unauthorized:
			s.SetFlash(w, Flash{Failure: "You can't delete the link!"})
			http.Redirect(w, r, "/", http.StatusFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	s.SetFlash(w, Flash{Success: "Successfully deleted link!"})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) executeUpdateLinkFormTemplate(w http.ResponseWriter, r *http.Request, f Flash, l *model.Link, aa []*model.Auth) {
	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"links/edit.css",
	}, nil, "links/edit.html", map[string]interface{}{
		"CurrentUsername": s.currentUsername(r),
		"Flash":           f,
		"Link":            l,
		"Auths":           aa,
	})
}
