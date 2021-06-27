package pending_transactions

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/backend_gandola_soft/database"
	"github.com/julienschmidt/httprouter"
)

type PendingTransaction struct {
	Id int
	Type string
	Amount float32
	Description string
	CreatedAt string
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