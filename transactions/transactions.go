package transactions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"example.com/backend_gandola_soft/database"
	"example.com/backend_gandola_soft/types"
	"example.com/backend_gandola_soft/utils"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// time.Sleep(5 * time.Second)
	fmt.Fprintf(w, "Server working")
}

func GetTransactions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transactions := []types.TransactionWithBalance{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query("SELECT transactions_with_balances.id, transactions_with_balances.type, transactions_with_balances.amount, transactions_with_balances.description, transactions_with_balances.balance, transactions_with_balances.executed, transactions_with_balances.created_at, actors.id, actors.name FROM transactions_with_balances, actors WHERE transactions_with_balances.actor = actors.id ORDER BY transactions_with_balances.id DESC;")
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		transaction := types.TransactionWithBalance{}
		if err := rows.Scan(&transaction.Id, &transaction.Type, &transaction.Amount, &transaction.Description, &transaction.Balance, &transaction.Executed, &transaction.CreatedAt, &transaction.Actor.Id, &transaction.Actor.Name); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		transactions = append(transactions, transaction)
	}
	json_transactions, err := json.Marshal(transactions)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_transactions)
}

// TODO: check sql injection issue
func CreateTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transaction := types.TransactionWithBalance{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data recibida no corresponde con una transacción")
		return
	}
	if transaction.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el tipo de transacción")
		return
	}
	if transaction.Type != "input" && transaction.Type != "output" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El tipo de transacción solo puede ser del tipo 'input' o 'output'")
		return
	}
	if transaction.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El monto de la transacción es menor a cero (0)")
		return
	}
	if transaction.Amount > float32(types.MaxTransactionAmount) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El monto de la transacción exede el máximo permitido")
		return
	}
	if transaction.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción debe poseer una descripción")
		return
	}
	if transaction.Actor.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción debe poseer un actor")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var actorId int
	getActorIdQuery := fmt.Sprintf("SELECT id FROM actors WHERE id=%v", transaction.Actor.Id)
	actorIdRow, err := db.Query(getActorIdQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer actorIdRow.Close()
	for actorIdRow.Next() {
		if err := actorIdRow.Scan(&actorId); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if actorId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El actor especificado no existe")
		return
	}

	var lastBalance float32
	var newBalance float32
	getLastBalanceQuery := "SELECT balance FROM transactions_with_balances ORDER BY id desc LIMIT 1;"
	lastTransactionRow, err := db.Query(getLastBalanceQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer lastTransactionRow.Close()
	for lastTransactionRow.Next() {
		if err := lastTransactionRow.Scan(&lastBalance); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	if transaction.Type == "input" {
		newBalance = lastBalance + transaction.Amount
	} else if transaction.Type == "output" {
		newBalance = lastBalance - transaction.Amount
	}
	if newBalance < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Su transacción no pudo ser ejecutada porque genera un balance menor a cero (0)")
		return
	}
	if newBalance > float32(types.MaxBalanceAmount) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Su transacción no pudo ser ejecutada porque excede el balance máximo permitido")
		return
	}

	var insertedId int
	insertTransactionQuery := fmt.Sprintf("INSERT INTO transactions_with_balances(type, amount, description, balance, actor) VALUES ('%v', '%v', '%v', '%v', '%v') RETURNING id;", transaction.Type, transaction.Amount, transaction.Description, newBalance, transaction.Actor.Id)

	rowsInsertedId, err := db.Query(insertTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsInsertedId.Close()
	for rowsInsertedId.Next() {
		if err := rowsInsertedId.Scan(&insertedId); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	insertedTransaction := types.TransactionWithBalance{}
	retrieveTransactionQuery := fmt.Sprintf("SELECT transactions_with_balances.id, transactions_with_balances.type, transactions_with_balances.amount, transactions_with_balances.description, transactions_with_balances.balance, transactions_with_balances.executed, transactions_with_balances.created_at, actors.id, actors.name FROM transactions_with_balances, actors WHERE transactions_with_balances.actor = actors.id AND transactions_with_balances.id = '%v';", insertedId)
	rowsRetrievedTransaction, err := db.Query(retrieveTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsRetrievedTransaction.Close()
	for rowsRetrievedTransaction.Next() {
		err = rowsRetrievedTransaction.Scan(&insertedTransaction.Id, &insertedTransaction.Type, &insertedTransaction.Amount, &insertedTransaction.Description, &insertedTransaction.Balance, &insertedTransaction.Executed, &insertedTransaction.CreatedAt, &insertedTransaction.Actor.Id, &insertedTransaction.Actor.Name)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	response, err := json.Marshal(insertedTransaction)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func PatchTransaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestId := ps.ByName("id")
	transactionId, err := strconv.Atoi(requestId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	partialTransaction := types.TransactionWithBalance{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	err = json.Unmarshal(body, &partialTransaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data enviada no corresponde con una transacción parcial")
		return
	}
	if transactionId <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No puede modificar la transacción cero")
		return
	}
	if partialTransaction.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción debe poseer una descripión")
		return
	}
	var updatedId int
	updateQuery := fmt.Sprintf("UPDATE transactions_with_balances SET description='%v' WHERE id='%v' RETURNING id;", partialTransaction.Description, transactionId)

	db := database.ConnectDB()
	defer db.Close()
	rowsId, err := db.Query(updateQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsId.Close()
	for rowsId.Next() {
		if err := rowsId.Scan(&updatedId); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if updatedId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción con el id %v no existe", transactionId)
		return
	}

	modifiedTransaction := types.TransactionWithBalance{}
	retrieveTransactionQuery := fmt.Sprintf("SELECT transactions_with_balances.id, transactions_with_balances.type, transactions_with_balances.amount, transactions_with_balances.description, transactions_with_balances.balance, transactions_with_balances.executed, transactions_with_balances.created_at, actors.id, actors.name FROM transactions_with_balances, actors WHERE transactions_with_balances.actor = actors.id AND transactions_with_balances.id = '%v';", updatedId)
	rowsRetrieveTransactionQuery, err := db.Query(retrieveTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsRetrieveTransactionQuery.Close()
	for rowsRetrieveTransactionQuery.Next() {
		err = rowsRetrieveTransactionQuery.Scan(&modifiedTransaction.Id, &modifiedTransaction.Type, &modifiedTransaction.Amount, &modifiedTransaction.Description, &modifiedTransaction.Balance, &modifiedTransaction.Executed, &modifiedTransaction.CreatedAt, &modifiedTransaction.Actor.Id, &modifiedTransaction.Actor.Name)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	response, err := json.Marshal(modifiedTransaction)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func DeleteLastTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	deletedTransactionId := types.IdResponse{
		Id: -1,
	}
	db := database.ConnectDB()
	defer db.Close()
	query := "DELETE FROM transactions_with_balances WHERE id != 1 AND id in (SELECT id FROM transactions_with_balances ORDER BY id desc LIMIT 1) RETURNING id;"
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&deletedTransactionId.Id); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	rollBackIdQuery := "SELECT setval('transactions_with_balances_id_seq', (SELECT last_value from transactions_with_balances_id_seq) - 1);"
	_, err = db.Query(rollBackIdQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	if deletedTransactionId.Id == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No quedan más transacciones por eliminar")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(deletedTransactionId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}

func GetLastTransactionId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	lastTransactionId := types.IdResponse{
		Id: -1,
	}
	db := database.ConnectDB()
	defer db.Close()
	query := "SELECT id FROM transactions_with_balances ORDER BY id desc LIMIT 1;"
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&lastTransactionId.Id); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if lastTransactionId.Id == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No existen más transacciones")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(lastTransactionId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}

func UnexecuteLastTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	lastTransaction := types.TransactionWithBalance{}
	db := database.ConnectDB()
	defer db.Close()
	query := "SELECT * FROM transactions_with_balances ORDER BY id DESC LIMIT 1;"
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&lastTransaction.Id, &lastTransaction.Type, &lastTransaction.Amount, &lastTransaction.Description, &lastTransaction.Balance, &lastTransaction.Actor.Id, &lastTransaction.Executed, &lastTransaction.CreatedAt); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	var insertedPendingTransactionId int
	insertPendingTransactionQuery := fmt.Sprintf("INSERT INTO pending_transactions(type, amount, description, actor, created_at) VALUES ('%v', '%v', '%v', '%v', '%v') RETURNING id;", lastTransaction.Type, lastTransaction.Amount, lastTransaction.Description, lastTransaction.Actor.Id, lastTransaction.CreatedAt)
	rowsId, err := db.Query(insertPendingTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsId.Close()
	for rowsId.Next() {
		if err := rowsId.Scan(&insertedPendingTransactionId); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	newPendingTransaction := types.PendingTransaction{}
	retrievePendingTransactionQuery := fmt.Sprintf("SELECT pending_transactions.id, pending_transactions.type, pending_transactions.amount, pending_transactions.description, pending_transactions.created_at, actors.id, actors.name FROM pending_transactions, actors WHERE pending_transactions.actor = actors.id AND pending_transactions.id = '%v';", insertedPendingTransactionId)
	rowsRetrievedTransaction, err := db.Query(retrievePendingTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsRetrievedTransaction.Close()
	for rowsRetrievedTransaction.Next() {
		err = rowsRetrievedTransaction.Scan(&newPendingTransaction.Id, &newPendingTransaction.Type, &newPendingTransaction.Amount, &newPendingTransaction.Description, &newPendingTransaction.CreatedAt, &newPendingTransaction.Actor.Id, &newPendingTransaction.Actor.Name)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	deleteQuery := "DELETE FROM transactions_with_balances WHERE id != 1 AND id in (SELECT id FROM transactions_with_balances ORDER BY id desc LIMIT 1) RETURNING id;"
	_, err = db.Query(deleteQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	rollBackIdQuery := "SELECT setval('transactions_with_balances_id_seq', (SELECT last_value from transactions_with_balances_id_seq) - 1);"
	_, err = db.Query(rollBackIdQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	response, err := json.Marshal(newPendingTransaction)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
