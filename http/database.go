package http

import "github.com/ravernkoh/jabba/model"

// Database represents the database.
type Database interface {
	CreateUser(*model.User) error
	FetchUser(username string) (*model.User, error)

	// TODO: Change to CreateUserLink and CreateVisitorLink
	CreateLink(*model.Link, *model.Visitor) error
	FetchLinks(slugs []string) ([]*model.Link, error)
	FetchLink(slug string) (*model.Link, error)

	PutVisitor(*model.Visitor) error
	FetchVisitor(token string) (*model.Visitor, error)
}
