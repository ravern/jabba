package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ravern/jabba/auth"
	"github.com/ravern/jabba/errors"
	"github.com/ravern/jabba/http/middleware"
	"github.com/ravern/jabba/model"
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
	s.executeLoginFormTemplate(w, r, flash, "")
}

// Login attempts to log the user in.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)

	var (
		username = r.FormValue("username")
		password = r.FormValue("password")
	)

	var f Flash

	u, err := s.Database.GetUser(username)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("failed to find user")

		model.DummyCheckPassword()

		f.Failure = "Invalid username or password."
		s.executeLoginFormTemplate(w, r, f, username)
		return
	}

	if !u.Registered {
		logger.Warn("user not registered")
		f.Failure = "Invalid username or password."
		s.executeLoginFormTemplate(w, r, f, username)
		return
	}

	if err := u.CheckPassword(password); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("wrong password")

		f.Failure = "Invalid username or password."
		s.executeLoginFormTemplate(w, r, f, username)
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

	http.Redirect(w, r, "/", http.StatusFound)
}

// Logout logs the user out.
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	s.DeleteCookie(w, "user")
	http.Redirect(w, r, "/", http.StatusFound)
}

// CreateUserForm renders the user creation form.
func (s *Server) CreateUserForm(w http.ResponseWriter, r *http.Request) {
	flash, _ := s.Flash(w, r)
	s.executeCreateUserFormTemplate(w, r, flash, &model.User{})
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
		f.Failure = "Passwords didn't match!"
		s.executeCreateUserFormTemplate(w, r, f, &model.User{
			Username: username,
			Email:    email,
		})
		return
	}

	u, err := model.NewUser(username, email)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("failed to create user")

		f.Failure = "Could not create user."
		s.executeCreateUserFormTemplate(w, r, f, &model.User{
			Username: username,
			Email:    email,
		})
		return
	}

	if err := u.SetPassword(password, confirmPassword); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("failed to set password on user")

		switch err.(errors.Error).Type {
		case errors.NotMatched:
			f.Failure = "Passwords didn't match!"
		default:
			f.Failure = "Could not create user."
		}

		s.executeCreateUserFormTemplate(w, r, f, u)
		return
	}

	if user.Registered {
		u.LinkSlugs = []string{}
		err = s.Database.CreateUser(u)
	} else {
		u.LinkSlugs = user.LinkSlugs
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

		s.executeCreateUserFormTemplate(w, r, f, u)
		return
	}

	s.SetFlash(w, Flash{Success: "Successfully registered user!"})
	http.Redirect(w, r, "/login", http.StatusFound)
}

// UpdateUserForm renders the user update form.
func (s *Server) UpdateUserForm(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	user := s.User(r)

	if !user.Registered {
		logger.Warn("unauthorized update user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	flash, _ := s.Flash(w, r)
	s.executeUpdateUserFormTemplate(w, r, flash, user)
}

// UpdateUser attempts to update the user.
func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger := middleware.Logger(r)
	user := s.User(r)

	username := user.Username
	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")

	var (
		password        = r.FormValue("password")
		newPassword     = r.FormValue("new_password")
		confirmPassword = r.FormValue("confirm_password")
	)

	var f Flash

	if err := user.CheckPassword(password); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Warn("unauthorized user update")

		f.Failure = "Invalid password!"
		s.executeUpdateUserFormTemplate(w, r, f, user)
		return
	}

	if newPassword != "" {
		if err := user.SetPassword(newPassword, confirmPassword); err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Warn("could not set password on user")

			switch err.(errors.Error).Type {
			case errors.NotMatched:
				f.Failure = "Passwords didn't match!"
			default:
				f.Failure = "Could not update user."
			}

			s.executeUpdateUserFormTemplate(w, r, f, user)
			return
		}
	}

	if err := user.Validate(); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("failed user validation")

		f.Failure = "Could not update user."

		user.Username = username
		s.executeUpdateUserFormTemplate(w, r, f, user)
		return
	}

	if err := s.Database.UpdateUserUsername(username, user); err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to update user")

		switch err.(errors.Error).Type {
		case errors.AlreadyExists:
			f.Failure = "User already exists."
		default:
			f.Failure = "Could not update user."
		}

		user.Username = username
		s.executeUpdateUserFormTemplate(w, r, f, user)
		return
	}

	logger.WithFields(logrus.Fields{
		"username": user.Username,
	}).Info("updated user")

	s.SetFlash(w, Flash{Success: "Successfully updated user!"})
	http.Redirect(w, r, "/user/edit", http.StatusFound)
}

func (s *Server) currentUsername(r *http.Request) string {
	user, ok := r.Context().Value(keyUser).(*model.User)
	if !ok {
		return ""
	}
	if !user.Registered {
		return ""
	}
	return user.Username
}

func (s *Server) executeLoginFormTemplate(w http.ResponseWriter, r *http.Request, f Flash, username string) {
	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"login.css",
	}, nil, "login.html", map[string]interface{}{
		"CurrentUsername": s.currentUsername(r),
		"Flash":           f,
		"Username":        username,
	})
}

func (s *Server) executeCreateUserFormTemplate(w http.ResponseWriter, r *http.Request, f Flash, u *model.User) {
	s.executeUserFormTemplate(w, r, f, u, "Register", "REGISTER", "/users", false)
}

func (s *Server) executeUpdateUserFormTemplate(w http.ResponseWriter, r *http.Request, f Flash, u *model.User) {
	s.executeUserFormTemplate(w, r, f, u, "User", "UPDATE", "/user", true)
}

func (s *Server) executeUserFormTemplate(w http.ResponseWriter, r *http.Request, f Flash, u *model.User, title string, submit string, action string, update bool) {
	executeTemplate(w, r, "layout.html", []string{
		"nav.css",
		"users/form.css",
	}, nil, "users/form.html", map[string]interface{}{
		"CurrentUsername": s.currentUsername(r),
		"Flash":           f,
		"User":            u,
		"Title":           title,
		"Submit":          submit,
		"Action":          action,
		"Update":          update,
	})
}
