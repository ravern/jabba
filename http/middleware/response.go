package middleware

import "net/http"

// SetContentType sets the content type of the response.
func SetContentType(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			w.Header().Set("Content-Type", contentType)
		})
	}
}
