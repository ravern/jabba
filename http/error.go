package http

import (
	"net/http"

	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/sirupsen/logrus"
)

// notFound renders the 404 page.
func notFound(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, "layout.html", nil, "error.html", struct {
		Message string
	}{
		Message: "404 Not Found",
	}); err != nil {
		middleware.Logger(r).WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to execute template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// internalServerError renders the 500 page.
func internalServerError(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, "layout.html", nil, "error.html", struct {
		Message string
	}{
		Message: "500 Internal Server Error",
	}); err != nil {
		middleware.Logger(r).WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to execute template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
