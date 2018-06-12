package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// SetRequestID sets the request ID in the context.
var SetRequestID = middleware.RequestID

// SetLogger sets the logger in the context.
//
// This should be used in conjunction with SetRequestID.
func SetLogger(l logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l = l.WithFields(logrus.Fields{
				"request_id": r.Context().Value(middleware.RequestIDKey),
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
		l.WithFields(logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		}).Info("received request")

		next.ServeHTTP(&rw, r)

		l = l.WithFields(logrus.Fields{
			"status_code": rw.statusCode,
		})
		if rw.statusCode >= http.StatusInternalServerError {
			l.Error("responded to request")
		} else {
			l.Info("responded to request")
		}
	})
}
