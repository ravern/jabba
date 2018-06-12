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
	executeTemplate(w, r, "layout.html", nil, "login.html", struct {
		Flash Flash
	}{
		Flash: flash,
	})
}

// Login attempts to log the user in.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
}

// CreateUserForm renders the user creation form.
func (s *Server) CreateUserForm(w http.ResponseWriter, r *http.Request) {
	flash, _ := RetrieveFlash(w, r)
	executeTemplate(w, r, "layout.html", nil, "users/new.html", struct {
		Flash Flash
	}{
		Flash: flash,
	})
}

// CreateUser attempts to create the user.
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	user := s.User(r)

	var (
		username        = r.FormValue("username")
		email           = r.FormValue("email")
		password        = r.FormValue("password")
		confirmPassword = r.FormValue("confirm_password")
	)

	if password != confirmPassword {
		f := Flash{Failure: "Passwords didn't match"}
		f.Save(w)

		http.Redirect(w, r, "/users/new", http.StatusFound)
		return
	}

	var err error
	user, err = model.NewUser(username, email, password, user.LinkSlugs)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to create user")

		f := Flash{Failure: "Could not create user."}
		f.Save(w)

		http.Redirect(w, r, "/users/new", http.StatusFound)
		return
	}

	if err := s.Database.UpdateUserUsername(s.User(r).Username, user); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to update user")

		switch err.(errors.Error).Type {
		case errors.AlreadyExists:
			f := Flash{Failure: "Username already exists."}
			f.Save(w)
		default:
			f := Flash{Failure: "Could not create user."}
			f.Save(w)
		}

		http.Redirect(w, r, "/users/new", http.StatusFound)
		return
	}

	f := Flash{Success: "Successfully registered user!"}
	f.Save(w)

	http.Redirect(w, r, "/login", http.StatusFound)
}
