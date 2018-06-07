package model

import (
	"net/url"
	"time"
)

// Link represents a shortened link.
type Link struct {
	Slug    string
	Title   string
	URL     *url.URL
	Created time.Time
}
