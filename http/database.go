package http

import "github.com/ravernkoh/jabba/model"

// Database represents the database.
type Database interface {
	CreateUser(*model.User) error
	UpdateUser(*model.User) error
	UpdateUserUsername(username string, u *model.User) error
	GetUser(username string) (*model.User, error)

	CreateLink(*model.Link, *model.User) error
	IncrementLinkCount(*model.Link)
	UpdateLinkSlug(slug string, l *model.Link, u *model.User) error
	GetLinks(slugs []string) ([]*model.Link, error)
	GetLink(slug string) (*model.Link, error)
	DeleteLink(*model.Link, *model.User) error

	UpdateAuths([]*model.Auth, *model.Link) error
	GetAuths(ids []string) ([]*model.Auth, error)
}
