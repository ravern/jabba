package model

import (
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/ravernkoh/jabba/errors"
	uuid "github.com/satori/go.uuid"
)

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

// NewAuth creates a new auth, converting and validating the method and values.
func NewAuth(method string, values string) (*Auth, error) {
	m, err := newMethod(method)
	if err != nil {
		return nil, err
	}

	v := strings.Split(values, ",")
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}

	a := &Auth{
		ID:     uuid.NewV4().String(),
		Method: m,
		Values: v,
	}

	if err := a.Validate(); err != nil {
		return nil, err
	}

	return a, nil
}

// Validate validates the auth.
func (a *Auth) Validate() error {
	if _, err := govalidator.ValidateStruct(a); err != nil {
		return newValidationError("auth", err)
	}

	switch a.Method {
	case MethodPassword:
		if len(a.Values) != 1 {
			return errors.Error{
				Type:    errors.Invalid,
				Message: "auth: invalid",
			}
		}
	case MethodGoogle:
		for _, v := range a.Values {
			if !govalidator.IsEmail(v) {
				return errors.Error{
					Type:    errors.Invalid,
					Message: "auth: invalid",
				}
			}
		}
	}

	return nil
}

// newMethod creates a new method based on the given string.
func newMethod(method string) (Method, error) {
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
