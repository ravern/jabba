package http

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/packr"
	logrus "github.com/sirupsen/logrus"
)

// Server serves the website.
type Server struct {
	Port     string
	Assets   packr.Box
	Logger   logrus.FieldLogger
	Database Database
}

// Listen listens for requests, blocking until an error occurs.
func (s *Server) Listen() error {
	return http.ListenAndServe(s.Port, s.router())
}

// router mounts all routes and returns the router.
func (s *Server) router() chi.Router {
	r := chi.NewRouter()

	// Mount static assets
	FileServer(r, "/", s.Assets)

	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve static
// files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
