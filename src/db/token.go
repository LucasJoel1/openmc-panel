package db

import (
	"os"
	"github.com/golang-jwt/jwt/v5"
)

var secret string = os.Getenv("TOKEN_SECRET")

type Token struct {
	id int
	username string
	permissions uint64
}

func CreateToken(_id int, _username string, _permissions uint64) (string, error) {
	tokenInfo := Token{
		id: _id,
		username: _username,
		permissions: _permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          tokenInfo.id,
		"username":    tokenInfo.username,
		"permissions": tokenInfo.permissions,
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}