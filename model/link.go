package model

import (
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/html"
)

// Link represents a shortened link.
type Link struct {
	Slug    string    `json:"slug" valid:"-"`
	Title   string    `json:"title" valid:"-"`
	URL     string    `json:"url" valid:"url"`
	Created time.Time `json:"created" valid:"-"`
	Count   int       `json:"count" valid:"-"`
}

// NewLink creates a new link with the given URL.
//
// Title generation is attempted, by visiting the URL and extracting the <title>.
// The link will default to 'Untitled' if this does not succeed.
func NewLink(url string) (*Link, error) {
	// Try to generate title automatically
	title := "Untitled"
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
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

	l := &Link{
		Slug:    string(slug),
		Title:   title,
		URL:     url,
		Created: time.Now(),
	}

	// Validate
	if _, err := govalidator.ValidateStruct(l); err != nil {
		return nil, err
	}

	return l, nil
}
