package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/sirupsen/logrus"
)

// Server serves the website.
type Server struct {
	Port     string
	Hostname string
	Logger   logrus.FieldLogger
	Database Database
}

// Listen listens for requests, blocking until an error occurs.
func (s *Server) Listen() error {
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
		middleware.ErrorPage(http.StatusNotFound, notFound),
		middleware.ErrorPage(http.StatusInternalServerError, internalServerError),

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
		r.Get("/{slug}", s.Redirect)
	})

	// Mount link routes
	r.Group(func(r chi.Router) {
		r.Use(
			// Authentication
			s.SetUser,
		)

		r.Post("/links", s.CreateLink)
		r.Post("/links/{slug}/delete", s.DeleteLink)
	})

	// Mount user routes
	r.Group(func(r chi.Router) {
		r.Use(
			// Authentication
			s.SetUser,
		)

		r.Get("/login", s.LoginForm)
		r.Post("/login", s.Login)

		r.Get("/users/new", s.CreateUserForm)
		r.Post("/users", s.CreateUser)
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
