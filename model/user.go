package model

import "time"

// User represents a user.
type User struct {
	Username string
	Email    string
	Password string
	Joined   time.Time
}
