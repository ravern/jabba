package http

import (
	"context"
	"net/http"

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
	// This is a lie!
	s.Logger.Infof("http: server listening on %s", s.Port)

	err := http.ListenAndServe(s.Port, s.Router())
	s.Logger.Error("http: server quit")
	return err
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

	// Mount main routes
	r.Get("/", s.Root)

	// Mount assets
	fileServer(r, "/public", assets)

	// Override not found handler to prevent "404 page not found" from
	// being sent in the response.
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
