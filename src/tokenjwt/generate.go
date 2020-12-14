package tokenjwt

import (
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Generate(username string, secret string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Username": username,
		"Expires":  time.Now().Add(ttl).Format(time.RFC3339),
	})

	tokenSigned, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return url.QueryEscape(tokenSigned), nil
}
