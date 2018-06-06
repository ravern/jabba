package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/sirupsen/logrus"
)

// Server serves the website.
type Server struct {
	Port     string
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
		middleware.SetLogger(s.Logger),
		middleware.LogRequest,

		// Error pages
		middleware.ErrorPage(http.StatusNotFound, notFound),
		middleware.ErrorPage(http.StatusInternalServerError, internalServerError),
	)

	// Mount static routes
	r.Get("/", s.Landing)

	// Mount assets
	fileServer(r, "/", assets)

	// Prevent any text from being rendered during a 404
	r.NotFound(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, middleware.Neuter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})))
}
