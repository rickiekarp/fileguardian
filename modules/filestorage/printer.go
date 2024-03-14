package filestorage

import (
	"database/sql"
	"fmt"
	"log"
)

func DisplayEntries(db *sql.DB) {
	row, err := db.Query("SELECT * FROM files ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		var id int
		var tag string
		var src string
		var dst string
		row.Scan(&id, &tag, &src, &dst)
		log.Println("File: ", id, " ", tag, " ", src, " ", dst)
	}
}

func PrintData(files []File) {
	for _, elem := range files {
		fmt.Println(elem)
	}
}
