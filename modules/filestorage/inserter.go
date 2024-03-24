package filestorage

import (
	"database/sql"
	"log"
)

func AddToStorage(db *sql.DB, files []File) {
	log.Println("Inserting file entries...")

	// create slice of files to add to the storage
	for _, file := range files {
		InsertEntry(db, file.Tag, file.Src, file.Dst)
	}

	log.Println("Done!")
}

func InsertEntry(db *sql.DB, tag string, src string, dst string) {
	insertSQL := `INSERT INTO files(tag, src, dst) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(tag, src, dst)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
