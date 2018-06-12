package http

import (
	"net/http"
	"time"
)

// SetCookie sets a cookie in the response.
func (s *Server) SetCookie(w http.ResponseWriter, key string, value interface{}) error {
	v, err := s.cookie.Encode(key, value)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  key,
		Value: v,
		Path:  "/",
	})

	return nil
}

// DeleteCookie deletes the cookie in the response.
func (s *Server) DeleteCookie(w http.ResponseWriter, key string) {
	http.SetCookie(w, &http.Cookie{
		Name:    key,
		Expires: time.Unix(0, 0),
		Path:    "/",
	})
}

// Cookie returns the cookie from the request.
func (s *Server) Cookie(r *http.Request, key string, value interface{}) error {
	c, err := r.Cookie(key)
	if err != nil {
		return err
	}
	return s.cookie.Decode(key, c.Value, value)
}
