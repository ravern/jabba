package model

import (
	"math/rand"
	"time"
)

// runes are the possible runes for generating slugs.
var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Link represents a shortened link.
type Link struct {
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
}

// NewLink creates a new link with the given URL.
//
// TODO: Add checks for URL validity
// TODO: Add automatic title generation
func NewLink(url string) *Link {
	// Generate the slug
	slug := make([]rune, 6)
	for i := range slug {
		slug[i] = runes[rand.Intn(len(runes))]
	}

	return &Link{
		Slug:    string(slug),
		Title:   "Untitled",
		URL:     url,
		Created: time.Now(),
	}
}
