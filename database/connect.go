package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB{
	connStr := "user=postgres password=1234 host=localhost port=5432 dbname=gandola_soft"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		return db
}