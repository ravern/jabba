package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Visitor represents an unregistered user.
type Visitor struct {
	Token     string    `json:"token"`
	Joined    time.Time `json:"joined"`
	LinkSlugs []string  `json:"link_slugs"`
}

// NewVisitor creates a new visitor with a random unique token.
func NewVisitor() *Visitor {
	return &Visitor{
		Token:     uuid.NewV4().String(),
		Joined:    time.Now(),
		LinkSlugs: []string{},
	}
}
