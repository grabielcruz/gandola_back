package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Transaction struct {
	Id       int
	Type     string
	Amount   float32
	Executed string
}

type Response struct {
	Message string
	Payload interface{}
}

var db *sql.DB

func init() {
	connStr := "user=postgres password=postgres host=localhost port=5432 dbname=gandola_soft"
	dbconn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connected to database")
	db = dbconn
}

func main() {
	defer db.Close()
	router := httprouter.New()
	router.GET("/", Index)

	router.GET("/transactions", GetTransactions)
	router.POST("/transactions", CreateTransaction)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Server working")
}

func GetTransactions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transactions := []Transaction{}
	rows, err := db.Query("SELECT * FROM transactions;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		transaction := Transaction{}
		if err := rows.Scan(&transaction.Id, &transaction.Type, &transaction.Amount, &transaction.Executed); err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, transaction)
	}
	json_transactions, err := json.Marshal(transactions)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_transactions)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transaction := Transaction{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
	}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		fmt.Fprintf(w, "La data enviada no corresponde con una transacción")
	}
	query := fmt.Sprintf("INSERT INTO transactions(type, amount, executed) VALUES ('%v', '%v', '%v') RETURNING id, type, amount, executed;", transaction.Type, transaction.Amount, transaction.Executed)
	inserted_transaction := Transaction{}
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&inserted_transaction.Id, &inserted_transaction.Type, &inserted_transaction.Amount, &inserted_transaction.Executed); err != nil {
			log.Fatal(err)
		}
	}
	response_data := Response{
		Message: "Transacción creada existosamente",
		Payload: inserted_transaction,
	}
	response, err := json.Marshal(response_data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
