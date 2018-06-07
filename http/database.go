package http

import "github.com/ravernkoh/jabba/model"

// Database represents the database.
type Database interface {
	CreateUser(*model.User) error
	FetchUser(username string) (*model.User, error)

	CreateLink(*model.Link) error
	FetchLinks(*model.User) ([]*model.Link, error)
	FetchLink(slug string) (*model.Link, error)
}
