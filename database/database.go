package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func createTable() {

	query := `
	CREATE TABLE IF NOT EXISTS users (
		chat_id INTEGER PRIMARY KEY,
		city TEXT
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Users table ready!")

}

func InitDatabase() {

	var err error

	DB, err = sql.Open("sqlite3", "bot.db")
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected!")

	createTable()

}
