package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ravernkoh/jabba/auth"
	"github.com/ravernkoh/jabba/errors"
	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/ravernkoh/jabba/model"
	"github.com/sirupsen/logrus"
)

// SetUser sets the user in the context, creating an anonymous one if not found.
//
// Must be placed after SetLogger.
func (s *Server) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.Logger(r)

		var u *model.User

		// Try to find in cookie
		var token string
		err := s.Cookie(r, "user", &token)
		if err == nil {
			// Decode the token
			var username string
			username, err = auth.ValidateToken(token, s.AuthSecret)
			if err == nil {
				logger.WithFields(logrus.Fields{
					"token": fmt.Sprintf("%s...", token[:10]),
				}).Info("decoded token")

				// Try to get from database
				u, err = s.Database.GetUser(username)
				if err == nil {
					logger.WithFields(logrus.Fields{
						"username": username,
					}).Info("got user")

					// Update the last visit
					u.LastVisit = time.Now()
					if err = s.Database.UpdateUser(u); err != nil {
						logger.WithFields(logrus.Fields{
							"username": u.Username,
						}).Error("failed to update user")
					}
				} else {
					logger.WithFields(logrus.Fields{
						"username": username,
						"err":      err,
					}).Error("failed to get user")
				}
			} else {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("failed to decode token")
			}
		} else {
			logger.Info("user cookie doesn't exist")
		}

		// Create new anonymous user, since either not found in cookie
		// or database.
		//
		// This is not in the else statement since err might change in
		// the previous if statement.
		if err != nil {
			u = model.NewAnonymousUser()

			if err := s.Database.CreateUser(u); err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("failed to create user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			token, err := auth.GenerateToken(u.Username, s.AuthSecret)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("failed to encode token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			s.SetCookie(w, "user", token)

			logger.Info("created new user")
		}

		ctx := context.WithValue(r.Context(), keyUser, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// User returns the user for the given request.
func (s *Server) User(r *http.Request) *model.User {
	return r.Context().Value(keyUser).(*model.User)
}

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

	u, err := model.NewUser(username, email, password, user.LinkSlugs)
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

	if user.Registered {
		err = s.Database.CreateUser(u)
	} else {
		err = s.Database.UpdateUserUsername(user.Username, u)
	}
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to update user")

		switch err.(errors.Error).Type {
		case errors.AlreadyExists:
			f.Failure = "Username already exists."
		default:
			f.Failure = "Could not create user."
		}

		executeCreateUserFormTemplate(w, r, f, u)
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
