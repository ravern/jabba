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
			ctx := context.WithValue(r.Context(), keyLogger, l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Logger returns the logger for the given request.
//
// Panics if the logger is not found.
func Logger(r *http.Request) logrus.FieldLogger {
	return r.Context().Value(keyLogger).(logrus.FieldLogger)
}

// LogRequest logs information about the request.
//
// It will panic if no logger is found in the context.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := responseWriter{w: w}
		l := Logger(r)
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
