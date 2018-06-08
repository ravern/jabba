package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Visitor represents an unregistered user.
type Visitor struct {
	Token     string    `json:"token"`
	LastVisit time.Time `json:"last_visit"`
	LinkSlugs []string  `json:"link_slugs"`
}

// NewVisitor creates a new visitor with a random unique token.
func NewVisitor() *Visitor {
	return &Visitor{
		Token:     uuid.NewV4().String(),
		LastVisit: time.Now(),
		LinkSlugs: []string{},
	}
}

// FindLinkSlug searches for a slug that belongs to the visitor and returns its
// index and flag to indicate whether it exists.
func (v *Visitor) FindLinkSlug(slug string) (int, bool) {
	for i, s := range v.LinkSlugs {
		if s == slug {
			return i, true
		}
	}
	return 0, false
}
