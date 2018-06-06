// Package middleware defines some generic middleware.
package middleware

// Key is used within contexts as a key.
type Key string

// Keys used by the middleware defined in this package.
const (
	KeyLogger Key = "logger"
)
