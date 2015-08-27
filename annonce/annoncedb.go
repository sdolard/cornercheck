package annonce

import (
	"database/sql"
	"log"
)

const (
	tableName  = "annonces"
	timeLayout = "02 Jan 06 15:04"
)

// TODO: select * from annonces where date(time)="2015-05-07";

// CreateAnnonceTable relative to db
func CreateAnnonceTable(db *sql.DB) {
	if db == nil {
		return
	}
	sqlStmt := `
		create table if not exists ` + tableName + ` (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			lbcID			TEXT,
			hRef			TEXT,
			title			TEXT,
			time			TEXT,
			category		TEXT,
			maxPrice		INT,
			minPrice		INT,
			town			TEXT,
			area			TEXT,
			placementString	TEXT
		);
		`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
		return
	}
	log.Printf("Table %v created", tableName)
}

// Insert ...
func Insert(cAnnonces chan []Annonce, db *sql.DB, channelCalls int) int {
	var annonces []Annonce
	for i := 0; i < channelCalls; i++ {
		annonces = append(annonces, <-cAnnonces...)
	}
	annCount := len(annonces)
	if len(annonces) > 0 {
		insert(annonces, db)
	}
	return annCount
}

// Insert an annonce in db
func insert(annonces []Annonce, db *sql.DB) {
	sqlStmt := `
		insert into ` + tableName + `(
			lbcID,
			hRef,
			title,
			time,
			category,
			maxPrice,
			minPrice,
			town,
			area,
			placementString
		)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
	transaction, err := db.Begin()
	if err != nil {
		log.Fatalf("Begin transation: %v\n", err)
		return
	}

	stmt, err := transaction.Prepare(sqlStmt)
	if err != nil {
		log.Fatalf("Transaction Prepare %q: %s\n", err, sqlStmt)
		return
	}

	log.Printf("Inserting %v annonces...", len(annonces))
	for _, ann := range annonces {
		// log.Printf("%v# %v: %v, %v-%v (%v), %v, %v, %v\n",
		// 	ann.Time.Format(timeLayout),
		// 	ann.Category,
		// 	ann.Title,
		// 	ann.MinPrice,
		// 	ann.MaxPrice,
		// 	ann.PriceString,
		// 	ann.PlacementString,
		// 	ann.LbcID(),
		// 	ann.HRef)

		_, err = stmt.Exec(
			ann.LbcID(),
			ann.HRef,
			ann.Title,
			ann.Time,
			ann.Category,
			ann.MaxPrice,
			ann.MinPrice,
			ann.Town,
			ann.Area,
			ann.PlacementString)
		if err != nil {
			log.Fatalf("stmt.Exec: %v. Try again", err)
		}
	}
	transaction.Commit()
	log.Printf("Inserting %v annonces: DONE", len(annonces))
}
