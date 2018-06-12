package http

import (
	"net/http"
	"time"
)

// Flash is the flash stored in the cookie.
type Flash struct {
	Success string
	Failure string
}

// Save sets the flash in the cookie to be used on the next request.
func (f Flash) Save(w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Name:  "flash-success",
		Value: f.Success,
		Path:  "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "flash-failure",
		Value: f.Failure,
		Path:  "/",
	})

	return nil
}

// RetrieveFlash retrieves the flash saved in the cookie and removes it.
func RetrieveFlash(w http.ResponseWriter, r *http.Request) (Flash, error) {
	var f Flash

	c, err := r.Cookie("flash-success")
	if err != nil {
		return f, err
	}
	f.Success = c.Value
	http.SetCookie(w, &http.Cookie{
		Name:    "flash-success",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})

	c, err = r.Cookie("flash-failure")
	if err != nil {
		return f, err
	}
	f.Failure = c.Value
	http.SetCookie(w, &http.Cookie{
		Name:    "flash-failure",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})

	return f, nil
}
