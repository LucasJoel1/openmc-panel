package db

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	_ "modernc.org/sqlite"
)

func SetupDB() error {
	db, err := OpenDB()

	if err != nil {
		return err
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	createUserTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			permissions BLOB NOT NULL,
			apiToken TEXT,
			admin_account BOOLEAN NOT NULL
		);
	`

	_, err = db.Exec(createUserTable)

	if err != nil {
		return err
	}

	generatedPassword, err := randomLowercaseString(30)

	if err != nil {
		return err
	}

	err = CreateUser("admin", string(generatedPassword), true)

	if err == nil {
		fmt.Printf("first time admin password generated: \033[1m%s\033[0m\nPlease change this on first login.", string(generatedPassword))
	}

	return nil
}

func randomLowercaseString(length int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(26))
		if err != nil {
			return "", err
		}
		b[i] = byte(97 + n.Int64()) // 97 is 'a'
	}
	return string(b), nil
}
