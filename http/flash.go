package http

import (
	"net/http"
)

// Flash is the flash stored in the cookie.
type Flash struct {
	Success string
	Failure string
}

// SetFlash sets the flash in the cookie to be used on the next request.
func (s *Server) SetFlash(w http.ResponseWriter, f Flash) error {
	return s.SetCookie(w, "flash", f)
}

// Flash retrieves the flash saved in the cookie and removes it.
func (s *Server) Flash(w http.ResponseWriter, r *http.Request) (Flash, error) {
	var f Flash

	if err := s.Cookie(r, "flash", &f); err != nil {
		return Flash{}, err
	}
	s.DeleteCookie(w, "flash")

	return f, nil
}
