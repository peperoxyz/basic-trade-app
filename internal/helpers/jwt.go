package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "your-256-bit-secret"

func GenerateToken(id uint, uuid string, email string) string {
	claims := jwt.MapClaims {
		"id": id,
		"uuid": uuid,
		"email": email,
		"exp": time.Now().Add(time.Minute * 10),
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := parseToken.SignedString([]byte(secretKey))
	if err != nil {
		return err.Error()
	}

	return signedToken
}