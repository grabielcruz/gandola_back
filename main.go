package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Transaction struct {
	id int
	amount float32
	executed string
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transactions := []Transaction{}
	connStr := "user=postgres password=postgres host=localhost port=5432 dbname=gandola_soft"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Print("Connected to database")

	rows, err := db.Query("SELECT * FROM transactions;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		transaction := Transaction{}
		if err := rows.Scan(&transaction.id, &transaction.amount, &transaction.executed); err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, transaction)
		// log.Printf("id %d amount %f executed %v", id, amount, executed)
	}
	
	fmt.Fprintf(w, "Welcome!\n %v", transactions)
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	log.Fatal(http.ListenAndServe(":8080", router))
}