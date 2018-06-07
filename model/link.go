package model

import (
	"time"
)

// Link represents a shortened link.
type Link struct {
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
}
