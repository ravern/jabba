package http

import "github.com/ravernkoh/jabba/model"

// Database represents the database.
type Database interface {
	CreateUser(*model.User) error
	UpdateUser(*model.User) error
	GetUser(username string) (*model.User, error)

	CreateLink(*model.Link, *model.User) error
	DeleteLink(slug string, u *model.User) error
	IncrementLinkCount(*model.Link)
	GetLinks(slugs []string) ([]*model.Link, error)
	GetLink(slug string) (*model.Link, error)
}
