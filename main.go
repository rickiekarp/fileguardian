package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const extension = "fg"
const storageName = "storage"

func main() {
	os.Remove(storageName + "." + extension)

	log.Println("Creating database...")
	file, err := os.Create(storageName + "." + extension)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("Database created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./"+storageName+"."+extension)
	defer sqliteDatabase.Close()
	createTable(sqliteDatabase)

	log.Println("Inserting file entries...")
	insertEntry(sqliteDatabase, "TagA", "FileA", "ObfuscatedA")
	insertEntry(sqliteDatabase, "TagA", "FileB", "ObfuscatedB")
	insertEntry(sqliteDatabase, "TagA", "FileC", "ObfuscatedC")
	insertEntry(sqliteDatabase, "TagA", "FileD", "ObfuscatedD")

	displayEntries(sqliteDatabase)
}

func createTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE files (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        "tag" varchar(100),
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

func insertEntry(db *sql.DB, tag string, src string, dst string) {
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

func displayEntries(db *sql.DB) {
	row, err := db.Query("SELECT * FROM files ORDER BY src")
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
