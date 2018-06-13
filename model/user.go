package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ravernkoh/jabba/errors"
	"golang.org/x/crypto/bcrypt"
)

// User represents a registered user.
type User struct {
	Username   string    `json:"username" valid:"stringlength(2|20)"`
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
func NewUser(username string, email string, password string) (*User, error) {
	// Check password length manually since after hashing its all the same
	if len(password) < 4 {
		return nil, errors.Error{
			Type:    errors.Invalid,
			Message: "user: invalid",
		}
	}

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
	}

	if err := u.Validate(); err != nil {
		return nil, newValidationError("user", err)
	}

	return u, nil
}

// Validate validates the user.
func (u *User) Validate() error {
	_, err := govalidator.ValidateStruct(u)
	return err
}

// CheckPassword checks whether the given password is correct.
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
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
