package middleware

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// SetLogger sets the logger in the context.
func SetLogger(l logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l = l.WithFields(logrus.Fields{
				"request_id": uuid.NewV4(),
				"path":       r.URL.Path,
			})
			ctx := context.WithValue(r.Context(), KeyLogger, l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// LogRequest logs information about the request.
//
// It will panic if no logger is found in the context.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responseWriter{w: w}
		l := r.Context().Value(KeyLogger).(logrus.FieldLogger)
		l.Info("http: received request")

		next.ServeHTTP(&rw, r)

		l = l.WithFields(logrus.Fields{
			"status_code": rw.statusCode,
		})
		if rw.statusCode >= http.StatusInternalServerError {
			l.Error("http: responded to request")
		} else {
			l.Info("http: responded to request")
		}
	})
}
