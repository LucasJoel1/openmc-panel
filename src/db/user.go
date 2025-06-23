package db

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
)

var usernameExpr, _ = regexp.Compile(`^[a-zA-Z0-9_.-]{3,30}$`)
var passwordExpr, _ = regexp.Compile(`^[a-zA-Z0-9!"#$%&'()*+,-.\/:;<=>?@[\]^_{|}~]{8,64}$`)

func generateSalt() (string, error) {
	var salt = make([]byte, 32)
	_, err := rand.Read(salt)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(salt), nil
}

func hashPassword(password, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt)) // Concatenate password and salt
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateUser(username string, password string, isAdmin bool) (err error) {
	db, err := OpenDB()

	if err != nil {
		return err
	}

	if !usernameExpr.Match([]byte(username)) {
		return fmt.Errorf("username invalid")
	}

	if !passwordExpr.Match([]byte(password)) {
		return fmt.Errorf("password invalid")
	}

	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username=?)`, username).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("user already exists")
	}

	salt, err := generateSalt()

	if err != nil {
		return fmt.Errorf("error generating salt")
	}

	hashedPassword := hashPassword(password, salt)

	permissions := make([]byte, 8)
	
	if isAdmin {
		binary.BigEndian.PutUint64(permissions, ^uint64(0))
	} else {
		binary.BigEndian.PutUint64(permissions, uint64(0))
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash, salt, permissions) VALUES (?, ?, ?, ?)", username, hashedPassword, salt, permissions)

	if err != nil {
		return err
	}

	return nil
}

func VerifyUserCreds(username string, password string) (jwt string, err error) {
	db, err := OpenDB()

	if err != nil {
		return "", err
	}

	var id int
	var passwordHash string
	var salt string
	var permissions []byte

	err = db.QueryRow("SELECT id, password_hash, salt, permissions FROM users WHERE username=?", username).Scan(&id, &passwordHash, &salt, &permissions)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}

	hashedCheckPassword := hashPassword(password, salt)

	parsedPerms := binary.BigEndian.Uint64(permissions)

	if subtle.ConstantTimeCompare([]byte(hashedCheckPassword), []byte(passwordHash)) != 1 {
		return "", fmt.Errorf("wrong password")
	}

	token, err := CreateToken(id, username, parsedPerms)

	if err != nil {
		return "", err
	}

	return token, nil
}