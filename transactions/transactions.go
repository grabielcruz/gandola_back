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

type TransactionWithBalance struct {
	Id          int
	Type        string
	Amount      float32
	Description string
	Executed    string
	Balance     float32
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
	transactions := []TransactionWithBalance{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM transactions_with_balances ORDER BY id ASC;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		transaction := TransactionWithBalance{}
		if err := rows.Scan(&transaction.Id, &transaction.Type, &transaction.Amount, &transaction.Description, &transaction.Balance, &transaction.Executed); err != nil {
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

	db := database.ConnectDB()
	defer db.Close()

	var lastBalance float32
	var newBalance float32
	getLastBalanceQuery := "SELECT balance FROM transactions_with_balances ORDER BY id desc LIMIT 1"
	lastTransactionRow, err := db.Query(getLastBalanceQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer lastTransactionRow.Close()
	for lastTransactionRow.Next() {
		if err := lastTransactionRow.Scan(&lastBalance); err != nil {
			log.Fatal(err)
		}
	}
	if transaction.Type == "input" {
		newBalance = lastBalance + transaction.Amount
	} else if transaction.Type == "output" {
		newBalance = lastBalance - transaction.Amount
		if (newBalance < 0) {
			response_data := Response{
				Message: "No se puede ejecutar esta transacción porque su balance no puede bajar de cero (0)",
				Payload: transaction,
			}
			response, err := json.Marshal(response_data)
			if err != nil {
				log.Fatal(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
			return
		}
	}

	insertedTransaction := TransactionWithBalance{}
	insertedTransactionQuery := fmt.Sprintf("INSERT INTO transactions_with_balances(type, amount, description, balance) VALUES ('%v', '%v', '%v', '%v') RETURNING id, type, amount, description, balance, executed;", transaction.Type, transaction.Amount, transaction.Description, newBalance)

	rows, err := db.Query(insertedTransactionQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&insertedTransaction.Id, &insertedTransaction.Type, &insertedTransaction.Amount, &insertedTransaction.Description, &insertedTransaction.Balance, &insertedTransaction.Executed); err != nil {
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

// TODO: forbid modification of transaction zero
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
	query := fmt.Sprintf("UPDATE transactions_with_balances SET description='%v' WHERE id='%v' RETURNING id, type, amount, description, balance, executed;", partialTransaction.Description, partialTransaction.Id)

	modifiedTransaction := TransactionWithBalance{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&modifiedTransaction.Id, &modifiedTransaction.Type, &modifiedTransaction.Amount, &modifiedTransaction.Description, &modifiedTransaction.Balance, &modifiedTransaction.Executed); err != nil {
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
	query := "DELETE FROM transactions_with_balances WHERE id != 1 AND id in (SELECT id FROM transactions_with_balances ORDER BY id desc LIMIT 1) RETURNING id;"
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
