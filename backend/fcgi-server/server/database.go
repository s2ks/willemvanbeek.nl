package server

import {
	"database/sql"
	"willemvanbeek.nl/backend/database"
}

func ActiveDB() (*sql.DB) {
	return database.Active()
}
