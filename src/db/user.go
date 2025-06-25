package db

import (
	"database/sql"
	"encoding/binary"
	"errors"
	"log"
	"regexp"
	"golang.org/x/crypto/bcrypt"
)

var usernameExpr, _ = regexp.Compile(`^[a-zA-Z0-9_.-]{3,30}$`)
var passwordExpr, _ = regexp.Compile(`^[a-zA-Z0-9!"#$%&'()*+,-.\/:;<=>?@[\]^_{|}~]{8,64}$`)

var ErrUserNotExist = errors.New("user not found")
var ErrWrongPassword = errors.New("wrong password")
var ErrUsernameInvalid = errors.New("username invalid")
var ErrPasswordInvalid = errors.New("password invalid")
var ErrUserExists = errors.New("user already exists")

type SmallUser struct {
	Username string `json:"username"`
	IsAdmin bool `json:"isAdmin"`
}

func CreateUser(username string, password string, isAdmin bool) (err error) {
	if !usernameExpr.Match([]byte(username)) {
		return ErrUsernameInvalid
	}

	if !passwordExpr.Match([]byte(password)) {
		return ErrPasswordInvalid
	}

	db, err := OpenDB()

	if err != nil {
		return err
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username=?)`, username).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return ErrUserExists
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

	_, err = db.Exec("INSERT INTO users (username, password_hash, permissions, admin_account) VALUES (?, ?, ?, ?)", username, hashedPassword, permissions, isAdmin)

	if err != nil {
		return err
	}

	return nil
}

func VerifyUserCreds(username string, password string) (jwt string, userPermissions uint64, err error) {
	db, err := OpenDB()

	if err != nil {
		return "", 0, err
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	var id int
	var passwordHash string
	var permissions []byte
	err = db.QueryRow("SELECT id, password_hash, permissions FROM users WHERE username=?", username).Scan(&id, &passwordHash, &permissions)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, ErrUserNotExist
		}
		return "", 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	if err != nil {
		return "", 0, ErrWrongPassword
	}

	parsedPerms := binary.BigEndian.Uint64(permissions)

	token, err := CreateToken(id, username, parsedPerms)

	if err != nil {
		return "", 0, err
	}

	return token, parsedPerms, nil
}

func GetUsersNames() (usernames []SmallUser, err error) {
	db, err := OpenDB()

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	rows, err := db.Query("SELECT username, admin_account FROM users")

	if err != nil {
		return nil, err
	}

	var users []SmallUser

	for rows.Next() {
		var username string
		var adminAccount bool
		err = rows.Scan(&username, &adminAccount)

		if err != nil {
			continue
		}

		users = append(users, SmallUser{Username: username, IsAdmin: adminAccount})
	}

	return users, nil
}