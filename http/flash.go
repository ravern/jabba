package http

import (
	"net/http"
	"time"
)

// setFlash sets the flash to be used on the next request.
func setFlash(w http.ResponseWriter, flash string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "flash",
		Value: flash,
	})
}

// flash returns the flash saved in the cookie, removing it in the process.
func flash(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie("flash")
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "flash",
		Expires: time.Unix(0, 0),
	})

	return c.Value, nil
}
