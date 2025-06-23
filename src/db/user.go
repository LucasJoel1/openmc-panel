package db

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"log"
	"regexp"
	"golang.org/x/crypto/bcrypt"
)

var usernameExpr, _ = regexp.Compile(`^[a-zA-Z0-9_.-]{3,30}$`)
var passwordExpr, _ = regexp.Compile(`^[a-zA-Z0-9!"#$%&'()*+,-.\/:;<=>?@[\]^_{|}~]{8,64}$`)

func CreateUser(username string, password string, isAdmin bool) (err error) {
	db, err := OpenDB()

	if err != nil {
		return err
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	permissions := make([]byte, 8)

	if isAdmin {
		binary.BigEndian.PutUint64(permissions, ^uint64(0))
	} else {
		binary.BigEndian.PutUint64(permissions, uint64(0))
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash, permissions) VALUES (?, ?, ?)", username, hashedPassword, permissions)

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
	var permissions []byte

	err = db.QueryRow("SELECT id, password_hash, permissions FROM users WHERE username=?", username).Scan(&id, &passwordHash, &permissions)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	if err != nil {
		return "", fmt.Errorf("wrong password")
	}

	parsedPerms := binary.BigEndian.Uint64(permissions)

	token, err := CreateToken(id, username, parsedPerms)

	if err != nil {
		return "", err
	}

	return token, nil
}
