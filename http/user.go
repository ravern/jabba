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
	flash, _ := s.Flash(w, r)
	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"login.css",
	}, nil, "login.html", map[string]interface{}{
		"Flash": flash,
	})
}

// Login attempts to log the user in.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
}

// CreateUserForm renders the user creation form.
func (s *Server) CreateUserForm(w http.ResponseWriter, r *http.Request) {
	flash, _ := s.Flash(w, r)
	executeCreateUserFormTemplate(w, r, flash, &model.User{})
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

	var f Flash

	if password != confirmPassword {
		f.Failure = "Passwords didn't match"
		executeCreateUserFormTemplate(w, r, f, &model.User{
			Username: username,
			Email:    email,
		})
		return
	}

	var err error
	user, err = model.NewUser(username, email, password, user.LinkSlugs)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("failed to create user")

		f.Failure = "Could not create user."
		executeCreateUserFormTemplate(w, r, f, &model.User{
			Username: username,
			Email:    email,
		})
		return
	}

	if err := s.Database.UpdateUserUsername(s.User(r).Username, user); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to update user")

		switch err.(errors.Error).Type {
		case errors.AlreadyExists:
			f.Failure = "Username already exists."
		default:
			f.Failure = "Could not create user."
		}

		executeCreateUserFormTemplate(w, r, f, user)
		return
	}

	s.SetFlash(w, Flash{Success: "Successfully registered user!"})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func executeCreateUserFormTemplate(w http.ResponseWriter, r *http.Request, f Flash, u *model.User) {
	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"users/new.css",
	}, nil, "users/new.html", map[string]interface{}{
		"Flash": f,
		"User":  u,
	})
}
