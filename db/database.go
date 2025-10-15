package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	createTable := `CREATE TABLE IF NOT EXISTS journeys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		start_time TEXT,
		end_time TEXT,
		journey_action TEXT,
		charge REAL,
		credit REAL,
		balance REAL,
		note TEXT,
		processed INTEGER DEFAULT 0,
		UNIQUE(date, start_time, end_time)
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
