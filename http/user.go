package http

import (
	"net/http"

	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/sirupsen/logrus"
)

// LoginForm renders the login form.
func (s *Server) LoginForm(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, "layout.html", nil, "login.html", nil); err != nil {
		middleware.Logger(r).WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to execute template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Login attempts to log the user in.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

}

// RegisterForm renders the registration form.
func (s *Server) RegisterForm(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, "layout.html", nil, "register.html", nil); err != nil {
		middleware.Logger(r).WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to execute template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Register attempts to register the user.
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, "layout.html", nil, "register.html", nil); err != nil {
		middleware.Logger(r).WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to execute template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
