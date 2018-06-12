package model

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

// Link represents a shortened link.
type Link struct {
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
	Count   int       `json:"count"`
}

// NewLink creates a new link with the given URL.
//
// Title generation is attempted, by visiting the URL and extracting the <title>.
// The link will default to 'Untitled' if this does not succeed.
//
// TODO: Add validations
func NewLink(rawURL string) (*Link, error) {
	// Validate URL
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return nil, fmt.Errorf("link: invalid URL: %v", err)
	}

	// Try to generate title automatically
	title := "Untitled"
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(rawURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		t := html.NewTokenizer(resp.Body)
		for {
			tt := t.Next()
			if tt != html.SelfClosingTagToken &&
				tt != html.StartTagToken &&
				tt != html.EndTagToken {
				continue
			}
			if t.Token().Data == "title" {
				t.Next()
				title = t.Token().Data
				break
			}
		}
	}

	// Generate the slug
	slug := generateToken(alphabet, 6)

	return &Link{
		Slug:    string(slug),
		Title:   title,
		URL:     rawURL,
		Created: time.Now(),
	}, nil
}
