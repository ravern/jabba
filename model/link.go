package model

import (
	"net/url"
	"time"
)

// Link represents a shortened link.
type Link struct {
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	URL     *url.URL  `json:"url"`
	Created time.Time `json:"created"`
}
