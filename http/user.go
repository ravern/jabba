package http

import (
	"net/http"

	"github.com/ravernkoh/jabba/errors"
	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/ravernkoh/jabba/model"
	"github.com/sirupsen/logrus"
)

// LoginForm renders the login form.
func (s *Server) LoginForm(w http.ResponseWriter, r *http.Request) {
	flash, _ := RetrieveFlash(w, r)
	if err := executeTemplate(w, "layout.html", nil, "login.html", struct {
		Flash Flash
	}{
		Flash: flash,
	}); err != nil {
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
	flash, _ := RetrieveFlash(w, r)
	if err := executeTemplate(w, "layout.html", nil, "register.html", struct {
		Flash Flash
	}{
		Flash: flash,
	}); err != nil {
		middleware.Logger(r).WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to execute template")

		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Register attempts to register the user.
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)

	var (
		username        = r.FormValue("username")
		email           = r.FormValue("email")
		password        = r.FormValue("password")
		confirmPassword = r.FormValue("confirm_password")
	)

	if password != confirmPassword {
		f := Flash{Failure: "Passwords didn't match"}
		f.Save(w)

		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	user, err := model.NewUser(username, email, password)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to create user")

		f := Flash{Failure: "Could not create user."}
		f.Save(w)

		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	if err := s.Database.CreateUser(user); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to create user")

		switch err.(errors.Error).Type {
		case errors.AlreadyExists:
			f := Flash{Failure: "Username already exists."}
			f.Save(w)
		default:
			f := Flash{Failure: "Could not create user."}
			f.Save(w)
		}

		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	f := Flash{Success: "Successfully registered user!"}
	f.Save(w)

	http.Redirect(w, r, "/login", http.StatusFound)
}
