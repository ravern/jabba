package middleware

import "net/http"

// ErrorPage runs the given function if the response contains the given status
// code.
func ErrorPage(statusCode int, f func(http.ResponseWriter)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := responseWriter{w: w}

			next.ServeHTTP(&rw, r)

			if rw.statusCode == statusCode {
				f(w)
			}
		})
	}
}
