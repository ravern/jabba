// Package middleware defines some generic middleware.
package middleware

import "net/http"

// Key is used within contexts as a key.
type Key string

// Keys used by the middleware defined in this package.
const (
	KeyLogger Key = "logger"
)

// responseWriter wraps a http.ResponseWriter, forwarding method calls while
// recording information about the response.
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
