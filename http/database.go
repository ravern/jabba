package http

import "github.com/ravernkoh/jabba/model"

// Database represents the database.
type Database interface {
	CreateUser(*model.User) error
	GetUser(username string) (*model.User, error)

	CreateUserLink(*model.Link, *model.User) error
	CreateVisitorLink(*model.Link, *model.Visitor) error
	DeleteUserLink(slug string, u *model.User) error
	DeleteVisitorLink(slug string, u *model.Visitor) error
	IncrementLinkCount(*model.Link)
	GetLinks(slugs []string) ([]*model.Link, error)
	GetLink(slug string) (*model.Link, error)

	PutVisitor(*model.Visitor) error
	GetVisitor(token string) (*model.Visitor, error)
}
