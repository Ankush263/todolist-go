package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("ABSKJOWUEO*!O!I&*@&*&@IJ@IOUEHKJ@BKBN@JBKJ@HO@U*#)*@)OUDWUWDHOIU")

func GenerateToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(token)
}

func ValidateToken(tokenstring string) (int64, error) {
	token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("Invalid Token")
	}

	claims := token.Claims.(jwt.MapClaims)
	return int64(claims["user_id"].(float64)), nil
}