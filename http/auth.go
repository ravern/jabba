package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Auth contains the values stored in the cookie for auth-related actions.
type Auth struct {
	Google string
}

// SetAuth sets the auth in the cookie to be used on the next request.
func (s *Server) SetAuth(w http.ResponseWriter, a Auth) error {
	return s.SetCookie(w, "auth", a)
}

// Auth retrieves the auth saved in the cookie and removes it.
func (s *Server) Auth(w http.ResponseWriter, r *http.Request) (Auth, error) {
	var a Auth

	if err := s.Cookie(r, "auth", &a); err != nil {
		return Auth{}, err
	}
	s.DeleteCookie(w, "auth")

	return a, nil
}

// AuthGoogle authenticates the user with Google and redirects back to the link.
func (s *Server) AuthGoogle(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)

	state, ok := r.URL.Query()["state"]
	if !ok || len(state) != 1 {
		logger.Error("failed to get state")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	code, ok := r.URL.Query()["code"]
	if !ok || len(code) != 1 {
		logger.Error("failed to get code")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := s.googleConfig.Exchange(oauth2.NoContext, code[0])
	if err != nil {
		logger.Error("failed to get token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := s.googleConfig.Client(oauth2.NoContext, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		logger.Error("failed to get userinfo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var j struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to parse json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.SetAuth(w, Auth{Google: j.Email}); err != nil {
		logger.Errorf("failed to set auth")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", state[0]), http.StatusFound)
}
