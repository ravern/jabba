package http

import (
	"context"
	"net/http"
	"time"

	"github.com/ravernkoh/jabba/http/middleware"
	"github.com/ravernkoh/jabba/model"
	"github.com/sirupsen/logrus"
)

// SetVisitor sets the visitor in the context, creating it if not found.
//
// Must be placed after SetLogger.
func (s *Server) SetVisitor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.Logger(r)

		var v *model.Visitor

		// Try to find in cookie
		c, err := r.Cookie("visitor")
		if err == nil {
			// Try to fetch from database
			v, err = s.Database.GetVisitor(c.Value)
			if err == nil {
				logger.WithFields(logrus.Fields{
					"token": c.Value,
				}).Info("fetched visitor")

				// Update the last visit
				v.LastVisit = time.Now()
				if err := s.Database.PutVisitor(v); err != nil {
					logger.WithFields(logrus.Fields{
						"token": v.Token,
					}).Warn("failed to put visitor")
				}
			} else {
				logger.WithFields(logrus.Fields{
					"token": c.Value,
				}).Warn("invalid visitor token found")
			}
		} else {
			logger.Info("failed to find visitor token in cookie")
		}

		// Create new visitor, since either not found in cookie or
		// database.
		//
		// This is not in the else statement since err might change in
		// the previous if statement.
		if err != nil {
			v = model.NewVisitor()

			if err := s.Database.PutVisitor(v); err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("failed to put visitor")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "visitor",
				Value: v.Token,
			})

			logger.Info("created new visitor")
		}

		ctx := context.WithValue(r.Context(), keyVisitor, v)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Visitor returns the visitor for the given request.
func (s *Server) Visitor(r *http.Request) *model.Visitor {
	return r.Context().Value(keyVisitor).(*model.Visitor)
}
