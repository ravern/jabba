package model

import "time"

// User represents a registered user.
type User struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Joined    time.Time `json:"joined"`
	LinkSlugs []string  `json:"link_slugs"`
}

// FindLinkSlug searches for a slug that belongs to the user and returns its
// index and flag to indicate whether it exists.
func (u *User) FindLinkSlug(slug string) (int, bool) {
	for i, s := range u.LinkSlugs {
		if s == slug {
			return i, true
		}
	}
	return 0, false
}
