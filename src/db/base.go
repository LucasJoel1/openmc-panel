package db

import "database/sql"

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./openmc.db")

	if err != nil {
		return nil, err
	}

	return db, nil
}