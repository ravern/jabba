package model

// Visitor represents an unregistered user.
type Visitor struct {
	Token     string   `json:"token"`
	LinkSlugs []string `json:"link_slugs"`
}
