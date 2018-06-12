// Package model defines all the models.
package model

import "math/rand"

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
