package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/securecookie"
	"github.com/ravern/jabba/http/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Server serves the website.
type Server struct {
	Port               string
	Hostname           string
	AuthSecret         string
	CookieHashKey      string
	CookieBlockKey     string
	GoogleClientID     string
	GoogleClientSecret string
	Logger             logrus.FieldLogger
	Database           Database

	cookie       *securecookie.SecureCookie
	googleConfig *oauth2.Config
}

// Listen listens for requests, blocking until an error occurs.
func (s *Server) Listen() error {
	s.cookie = securecookie.New([]byte(s.CookieHashKey), []byte(s.CookieBlockKey))

	s.googleConfig = &oauth2.Config{
		ClientID:     s.GoogleClientID,
		ClientSecret: s.GoogleClientSecret,
		RedirectURL:  s.Hostname + "/auth/google",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return http.ListenAndServe(s.Port, s.Router())
}

// Router mounts all routes on a router and returns it.
func (s *Server) Router() chi.Router {
	r := chi.NewRouter()

	// Global middleware
	r.Use(
		// Logging
		middleware.SetRequestID,
		middleware.SetLogger(s.Logger),
		middleware.LogRequest,

		// Error pages
		middleware.ErrorPage(http.StatusNotFound, s.notFound),
		middleware.ErrorPage(http.StatusInternalServerError, s.internalServerError),
		middleware.ErrorPage(http.StatusUnauthorized, s.unauthorized),

		// Response
		middleware.SetContentType("text/html; charset=utf-8"),
	)

	// Mount assets
	fileServer(r, "/public", assets)

	// Mount root routes
	r.Group(func(r chi.Router) {
		r.Use(
			// Authentication
			s.SetUser,
		)

		r.Get("/", s.Index)
		r.With(s.SetLink).Get("/{slug}", s.Redirect)
	})

	// Mount auth routes
	r.Group(func(r chi.Router) {
		r.Get("/auth/google", s.AuthGoogle)
	})

	// Mount link routes
	r.Group(func(r chi.Router) {
		r.Use(
			// Authentication
			s.SetUser,
		)

		r.Post("/links", s.CreateLink)
		r.Group(func(r chi.Router) {
			r.Use(
				s.SetLink,
				s.RequireLink,
			)
			r.Get("/links/{slug}/edit", s.UpdateLinkForm)
			r.Post("/links/{slug}", s.UpdateLink)
			r.Post("/links/{slug}/delete", s.DeleteLink)
		})
	})

	// Mount user routes
	r.Group(func(r chi.Router) {
		r.Use(
			// Authentication
			s.SetUser,
		)

		r.Get("/login", s.LoginForm)
		r.Post("/login", s.Login)
		r.Post("/logout", s.Logout)

		r.Get("/users/new", s.CreateUserForm)
		r.Post("/users", s.CreateUser)
		r.Get("/user/edit", s.UpdateUserForm)
		r.Post("/user", s.UpdateUser)
	})

	// Override not found handler to prevent "404 page not found" from
	// being sent in the response.
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}

// fileServer conveniently sets up a http.FileServer handler to serve assets
// from a http.FileSystem.
//
// https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go#L26-L44
//
// Added Neuter for directories.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.With(middleware.Neuter).Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
