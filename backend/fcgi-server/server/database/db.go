package database

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	File       string
	Connection *sql.DB
	Init       bool
}

var activeDB Database

func Open(dbfile string) (*sql.DB, error) {
	if activeDB.Init == true {
		return activeDB.Connection, nil
	}

	/* sql.Open creates the file if it doesn't exist */
	_, err := os.Stat(dbfile)

	if err != nil {
		return nil, err
	}

	if db, err := sql.Open("sqlite3", activeDB.File); err != nil {
		return nil, err
	} else {
		activeDB.Connection = db
		activeDB.Init = true
		activeDB.File = dbfile
		return db, nil
	}
}

func Active() (*sql.DB, error) {
	if activeDB.Connection == nil {
		return nil, fmt.Errorf("No active database")
	} else {
		return activeDB.Connection, nil
	}
}
