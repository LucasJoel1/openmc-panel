package db

import (
	"log"

	_ "modernc.org/sqlite"
)

func SetupDB() error {
	db, err := OpenDB()

	if err != nil {
		return err
	}

	defer db.Close()

	createUserTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			salt TEXT NOT NULL,
			permissions BLOB NOT NULL,
			apiToken TEXT
		);
	`

	_, err = db.Exec(createUserTable)

	if err != nil {
		return err
	}

	err = CreateUser("admin", "Password1", true)
	
	if err != nil {
		log.Println(err)
	}

	return nil
}