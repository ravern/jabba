package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ravernkoh/jabba/auth"
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

				// Try to fetch from database
				u, err = s.Database.GetUser(username)
				if err == nil {
					logger.WithFields(logrus.Fields{
						"username": username,
					}).Info("fetched user")

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
					}).Error("failed to fetch user")
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
