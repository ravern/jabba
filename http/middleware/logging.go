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

// responseWriter wraps a http.ResponseWriter, forwarding method calls while
// recording information about the response for logging.
type responseWriter struct {
	w          http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.w.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (n int, err error) {
	return rw.w.Write(b)
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
		if rw.statusCode >= 500 {
			l.Error("http: responded to request")
		} else {
			l.Info("http: responded to request")
		}
	})
}
