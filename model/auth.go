package model

import "fmt"

// Auth represents a form of authentication for the link.
type Auth struct {
	ID     string   `json:"id" valid:"-"`
	Method Method   `json:"method" valid:"-"`
	Values []string `json:"values" valid:"-"`
}

// Method is a method of authentication (e.g. Google, GitHub or password).
type Method int

// Supported methods of authentication.
const (
	MethodPassword Method = iota // password (stored in plaintext)
	MethodGoogle                 // email
	MethodJabba                  // username
)

// NewMethod creates a new method based on the given string.
func NewMethod(method string) (Method, error) {
	switch method {
	case "password":
		return MethodPassword, nil
	case "google":
		return MethodGoogle, nil
	case "jabba":
		return MethodJabba, nil
	}
	return 0, fmt.Errorf("auth: invalid method provided")
}
