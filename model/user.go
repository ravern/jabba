package model

import "time"

// User represents a user.
type User struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Joined    time.Time `json:"joined"`
	LinkSlugs []string  `json:"link_slugs"`
}
