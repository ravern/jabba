// Package model defines all the models.
package model

import (
	"fmt"
	"math/rand"

	"golang.org/x/crypto/bcrypt"

	"github.com/asaskevich/govalidator"
	"github.com/ravernkoh/jabba/errors"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// DummyCheckPassword is a fake password check to prevent timing attacks.
func DummyCheckPassword() {
	bcrypt.CompareHashAndPassword([]byte(""), []byte(""))
}

// alphabet is the alphabet.
var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generateToken generates a token of the given length, using the given runes.
func generateToken(runes []rune, length int) string {
	r := make([]rune, length)
	for i := range r {
		r[i] = runes[rand.Intn(len(runes))]
	}
	return string(r)
}

// newValidationError generates an error from the raw validation error.
func newValidationError(name string, err error) errors.Error {
	if errs, ok := err.(govalidator.Errors); ok {
		if _, ok := errs[0].(govalidator.Error); ok {
			return errors.Error{
				Type:    errors.Invalid,
				Message: fmt.Sprintf("%s: invalid", name),
			}
		}
	}
	return errors.Error{
		Type:    errors.Unknown,
		Message: fmt.Sprintf("%s: failed to create", name),
	}
}
