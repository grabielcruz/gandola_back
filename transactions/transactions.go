package transactions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"example.com/backend_gandola_soft/database"
	"github.com/julienschmidt/httprouter"
)

type Transaction struct {
	Id          int
	Type        string
	Amount      float32
	Description string
	Executed    string
}

type Response struct {
	Message string
	Payload interface{}
}

type PartialTransaction struct {
	Id          int
	Description string
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Server working")
}

func GetTransactions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transactions := []Transaction{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM transactions ORDER BY id ASC;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		transaction := Transaction{}
		if err := rows.Scan(&transaction.Id, &transaction.Type, &transaction.Amount, &transaction.Description, &transaction.Executed); err != nil {
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
	query := fmt.Sprintf("INSERT INTO transactions(type, amount, description) VALUES ('%v', '%v', '%v') RETURNING id, type, amount, description, executed;", transaction.Type, transaction.Amount, transaction.Description)
	insertedTransaction := Transaction{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&insertedTransaction.Id, &insertedTransaction.Type, &insertedTransaction.Amount, &insertedTransaction.Description, &insertedTransaction.Executed); err != nil {
			log.Fatal(err)
		}
	}
	response_data := Response{
		Message: "Transacción creada existosamente",
		Payload: insertedTransaction,
	}
	response, err := json.Marshal(response_data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func PatchTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	partialTransaction := PartialTransaction{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
	}
	err = json.Unmarshal(body, &partialTransaction)
	if err != nil {
		fmt.Fprintf(w, "La data enviada no corresponde con una transacción parcial")
	}
	query := fmt.Sprintf("UPDATE transactions SET description='%v' WHERE id='%v' RETURNING id, type, amount, description, executed;", partialTransaction.Description, partialTransaction.Id)

	modifiedTransaction := Transaction{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&modifiedTransaction.Id, &modifiedTransaction.Type, &modifiedTransaction.Amount, &modifiedTransaction.Description, &modifiedTransaction.Executed); err != nil {
			log.Fatal(err)
		}
	}

	if modifiedTransaction.Id == 0 {
		response_data := Response{
			Message: "La transacción con el id indicado no existe",
			Payload: partialTransaction.Id,
		}
		response, err := json.Marshal(response_data)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return
	}

	response_data := Response{
		Message: "Transacción modificada existosamente",
		Payload: modifiedTransaction,
	}
	response, err := json.Marshal(response_data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func DeleteLastTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	deletedTransactionId := -1
	db := database.ConnectDB()
	defer db.Close()
	query := "DELETE FROM transactions WHERE id in (SELECT id FROM transactions ORDER BY id desc LIMIT 1) RETURNING id;"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&deletedTransactionId); err != nil {
			log.Fatal(err)
		}
	}

	if deletedTransactionId == -1 {
		response_data := Response{
			Message: "No quedan más transacciones por eliminar por lo que no se pudo eliminar ninguna transacción",
			Payload: deletedTransactionId,
		}
		response, err := json.Marshal(response_data)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return
	}

	response_data := Response{
		Message: "Transacción eliminada exitosamente",
		Payload: deletedTransactionId,
	}
	response, err := json.Marshal(response_data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
