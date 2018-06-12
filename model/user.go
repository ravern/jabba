package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ravernkoh/jabba/errors"
	"golang.org/x/crypto/bcrypt"
)

// User represents a registered user.
type User struct {
	Username   string    `json:"username" valid:"-"`
	Registered bool      `json:"registered" valid:"-"`
	Email      string    `json:"email" valid:"email"`
	Password   string    `json:"password" valid:"-"`
	Joined     time.Time `json:"joined" valid:"-"`
	LastVisit  time.Time `json:"last_visit" valid:"-"`
	LinkSlugs  []string  `json:"link_slugs" valid:"-"`
}

// NewAnonymousUser creates a new anonymous user.
//
// Username is a random alphabetic string.
func NewAnonymousUser() *User {
	return &User{
		Username:  generateToken(alphabet, 8),
		LinkSlugs: []string{},
	}
}

// NewUser creates a new user.
//
// The given password will be hashed and stored in the password field.
func NewUser(username string, email string, password string, linkSlugs []string) (*User, error) {
	// Hash the passwords (use default cost)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), -1)
	if err != nil {
		return nil, errors.Error{
			Type:    errors.FailedHash,
			Message: "user: failed to hash password",
		}
	}

	u := &User{
		Username:   username,
		Registered: true,
		Email:      email,
		Password:   string(passwordHash),
		Joined:     time.Now(),
		LastVisit:  time.Now(),
		LinkSlugs:  linkSlugs,
	}

	// Validate
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return nil, newValidationError("user", err)
	}

	return u, nil
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
