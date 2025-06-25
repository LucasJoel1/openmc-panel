package db

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var secret []byte = getJWTSecret()

type Token struct {
	id          int
	username    string
	permissions uint64
	exp         int64
}

func getJWTSecret() []byte {
	godotenv.Load()
	secret := os.Getenv("TOKEN_SECRET")

	if secret == "" {
		panic("TOKEN_SECRET not set")
	}

	return []byte(secret)
}

func CreateToken(_id int, _username string, _permissions uint64) (string, error) {
	tokenInfo := Token{
		id:          _id,
		username:    _username,
		permissions: _permissions,
		exp:         time.Now().Add(24 * 30 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          tokenInfo.id,
		"username":    tokenInfo.username,
		"permissions": tokenInfo.permissions,
		"exp":         tokenInfo.exp, // 30 days
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
