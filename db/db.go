package db

import (
	"database/sql"
	"log"

	// sqldriver
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFilePath = "./lbc.db"
)

var db *sql.DB

// Open database
func Open() *sql.DB {
	if db != nil {
		return db
	}
	_db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Fatal(err)
		return db
	}
	db = _db
	log.Printf("Db opened: %v", dbFilePath)
	return db
}

// Close Database
func Close() {
	if db != nil {
		db.Close()
		db = nil
	}
}
