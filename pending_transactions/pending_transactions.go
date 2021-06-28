package pending_transactions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"example.com/backend_gandola_soft/database"
	"github.com/julienschmidt/httprouter"
)

type PendingTransaction struct {
	Id          int
	Type        string
	Amount      float32
	Description string
	CreatedAt   string
}

type PartialPendingTransaction struct {
	Id          int
	Type        string
	Amount      float32
	Description string
}

func GetPendingTransactions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transactions := []PendingTransaction{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM pending_transactions ORDER BY id ASC;")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		transaction := PendingTransaction{}
		err = rows.Scan(&transaction.Id, &transaction.Type, &transaction.Amount, &transaction.Description, &transaction.CreatedAt)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}
	json_transactions, err := json.Marshal(transactions)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_transactions)
}

func CreatePendingTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transaction := PartialPendingTransaction{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		fmt.Fprintf(w, "La data recibida no corresponde con una transacción pendiente")
	}
	if transaction.Type == "" || transaction.Type != "input" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el tipo de transacción pendiente")
		return
	}
	if transaction.Type != "input" && transaction.Type != "output" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El tipo de transacción solo puede ser del tipo 'input' o 'output'")
		return
	}
	if transaction.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El monto de la transacción pendiente debe ser mayor a cero")
		return
	}
	if transaction.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente debe poseer una descripción")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	insertedTransaction := PendingTransaction{}
	insertedTransactionQuery := fmt.Sprintf("INSERT INTO pending_transactions(type, amount, description) VALUES ('%v', '%v', '%v') RETURNING id, type, amount, description, created_at;", transaction.Type, transaction.Amount, transaction.Description)

	rows, err := db.Query(insertedTransactionQuery)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		err = rows.Scan(&insertedTransaction.Id, &insertedTransaction.Type, &insertedTransaction.Amount, &insertedTransaction.Description, &insertedTransaction.CreatedAt)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	response, err := json.Marshal(insertedTransaction)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func PatchPendingTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	newPendingTransaction := PendingTransaction{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
	}
	err = json.Unmarshal(body, &newPendingTransaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data enviada no corresponde con un transacción pendiente parcial")
		return
	}
	if newPendingTransaction.Id <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No puede modificar la transacción pendiente cero")
		return
	}
	if newPendingTransaction.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente debe poseer una descripión")
		return
	}
	if (newPendingTransaction.Type != "input" && newPendingTransaction.Type != "output") {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El tipo de la transacción debe ser 'input' o 'output'")
		return
	}
	if newPendingTransaction.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente debe poseer un monto mayor que cero")
		return
	}
	query := fmt.Sprintf("UPDATE pending_transactions SET type='%v', amount='%v', description='%v' WHERE id='%v' RETURNING id, type, amount, description, created_at;", newPendingTransaction.Type, newPendingTransaction.Amount, newPendingTransaction.Description, newPendingTransaction.Id)
	modifiedPendingTransaction := PendingTransaction{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&modifiedPendingTransaction.Id, &modifiedPendingTransaction.Type, &modifiedPendingTransaction.Amount, &modifiedPendingTransaction.Description, &modifiedPendingTransaction.CreatedAt)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if modifiedPendingTransaction.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente con el id %v no existe", newPendingTransaction.Id)
		return
	}
	response, err := json.Marshal(modifiedPendingTransaction)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}