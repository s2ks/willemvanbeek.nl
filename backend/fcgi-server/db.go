package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

type Database struct {
	File       string
	Connection *sql.DB
	Init       bool
}

var activeDB Database

func DatabaseInit(dbPath string) (*sql.DB, error) {
	var err error

	if activeDB.Init {
		return activeDB.Connection, nil
	}

	if _, err = os.Stat(dbPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dbPath, 0777); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if exec, err := os.Executable(); err != nil {
		return nil, err
	} else {
		activeDB.File = dbPath + filepath.Base(exec) + ".db"
	}

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
