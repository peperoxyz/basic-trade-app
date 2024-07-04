package helpers

import (
	"basic-trade-app/config"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = config.SecretKey

func GenerateToken(id uint, uuid string, email string) string {
	claims := jwt.MapClaims {
		"id": id,
		"uuid": uuid,
		"email": email,
		"exp": time.Now().Add(time.Minute * 72),
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := parseToken.SignedString([]byte(secretKey))
	if err != nil {
		return err.Error()
	}

	return signedToken
}

// function untuk ngecek apakah sebuah token berhak mengakses endpoint tertentu (middleware) 
func VerifyToken(ctx *gin.Context) (interface{}, error) {
	headerToken := ctx.Request.Header.Get("Authorization")

	// pengecekan apakah terdapat kata "Bearer" sebelum token
	bearer := strings.HasPrefix(headerToken, "Bearer")
	if !bearer {
		return nil, errors.New("You have to sign in to proceed")
	}

	stringToken := strings.Split(headerToken, " ")[1]

	token, _ := jwt.Parse(stringToken, func (t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Sign in to proceed")
		}
		return []byte(secretKey), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("Sign in to proceed")
	}

	expClaim, exists := claims["exp"]
	if !exists {
		return nil, errors.New("Expire claim is missing")
	}

	expStr, ok := expClaim.(string)
	if !ok {
		return nil, errors.New("Expire claim is not a valid type")
	}

	expTime, err := time.Parse(time.RFC3339, expStr)
	if err != nil {
		return nil, errors.New("Error parsing exp time")
	}

	if time.Now().After(expTime) {
		return nil, errors.New("Token is expired")
	}

	return token.Claims.(jwt.MapClaims), nil
}