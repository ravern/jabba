package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a registered user.
type User struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Joined    time.Time `json:"joined"`
	LinkSlugs []string  `json:"link_slugs"`
}

// NewUser creates a new user.
//
// The given password will be hashed and stored in the password field.
//
// TODO: Add validations
func NewUser(username string, email string, password string) (*User, error) {
	// Hash the passwords (use default cost)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), -1)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:  username,
		Email:     email,
		Password:  string(passwordHash),
		Joined:    time.Now(),
		LinkSlugs: []string{},
	}, nil
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
