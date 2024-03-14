package filestorage

import (
	"database/sql"
	"log"
)

func CreateTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE files (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        "tag" varchar(64),
		"src" TEXT,
		"dst" TEXT
	  );`

	log.Println("Create table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Table created")
}
