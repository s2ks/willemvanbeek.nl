package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	File       string
	Connection *sql.DB
	Init       bool
}

var activeDB Database

func DatabaseInit(dbPath string) (*sql.DB, error) {
	if activeDB.Init {
		return activeDB.Connection, nil
	}

	activeDB.File = dbPath

	/* sql.Open creates the file if it doesn't exist */
	if db, err := sql.Open("sqlite3", activeDB.File); err != nil {
		return nil, err
	} else {
		activeDB.Connection = db
		activeDB.Init = true
		return db, nil
	}
}

func ActiveDatabase() (*sql.DB, error) {
	if activeDB.Connection == nil {
		return nil, fmt.Errorf("No active database")
	} else {
		return activeDB.Connection, nil
	}
}

func ActiveDatabaseFile() (string, error) {
	if activeDB.Connection == nil {
		return "", fmt.Errorf("No active database")
	} else if activeDB.File == "" {
		return "", fmt.Errorf("No active database file")
	} else {
		return activeDB.File, nil
	}
}
