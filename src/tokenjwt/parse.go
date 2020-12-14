package tokenjwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Parse(tokenSigned string, secret string) (string, time.Time, error) {
	token, err := jwt.ParseWithClaims(strings.Replace(tokenSigned, "Bearer ", "", -1), &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC), err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		expires, _ := time.Parse(time.RFC3339, fmt.Sprintf("%s", (*claims)["Expires"]))
		return fmt.Sprintf("%s", (*claims)["Username"]), expires, nil
	}

	return "", time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC), fmt.Errorf("Failed token validation")
}
