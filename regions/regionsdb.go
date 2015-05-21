package regions

import (
	"database/sql"
	"log"
)

const (
	areasTableName   = "areas"
	regionsTableName = "regions"
	sqlAreasTableStr = `
		create table if not exists ` + areasTableName + ` (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name	TEXT
		);
		`
	sqlRegionsTableStr = `
			create table if not exists ` + regionsTableName + ` (
				id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
				name	TEXT,

			);
			`
)

// CreateTable relative to db
func CreateTable(db *sql.DB) {
	createTable(areasTableName, db)
	createTable(regionsTableName, db)
}

func createTable(tableName string, db *sql.DB) {
	if db == nil {
		return
	}
	sqlStmt := `
		create table if not exists ` + tableName + ` (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name	TEXT
		);
		`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
		return
	}
	log.Printf("Table %v created", tableName)
}
