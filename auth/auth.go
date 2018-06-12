package auth

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateToken generates a token with the username stored.
func GenerateToken(username string, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	c := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7),
	}
	token.Claims = c

	r, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return r, nil
}

// ValidateToken checks whether the token is valid and returns the username.
func ValidateToken(tokenString, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("auth: unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("auth: token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("auth: could not extract claims from token")
	}

	return claims["username"].(string), nil
}
