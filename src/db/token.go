package db

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var secret []byte = GetJWTSecret()

type Token struct {
	id          int
	username    string
	permissions string
	version     int
	exp         int64
}

func GetJWTSecret() []byte {
	godotenv.Load()
	secret := os.Getenv("TOKEN_SECRET")

	if secret == "" {
		panic("TOKEN_SECRET not set")
	}

	return []byte(secret)
}

func CreateToken(_id int, _username string, _permissions uint64, _tokenVersion int) (string, error) {
	tokenInfo := Token{
		id:          _id,
		username:    _username,
		permissions: strconv.FormatUint(_permissions, 10),
		version: _tokenVersion,
		exp: time.Now().Add(24 * 30 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":           tokenInfo.id,
		"username":     tokenInfo.username,
		"permissions":  tokenInfo.permissions,
		"tokenVersion": tokenInfo.version,
		"exp":          tokenInfo.exp, // 30 days
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
