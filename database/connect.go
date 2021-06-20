package database

import (
	"database/sql"
	"log"
)

func ConnectDB() *sql.DB{
	connStr := "user=postgres password=postgres host=localhost port=5432 dbname=gandola_soft"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Connected to database")
		return db
}