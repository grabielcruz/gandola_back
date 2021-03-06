package pending_transactions

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

func GetPendingTransactions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transactions := []types.PendingTransaction{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query("SELECT pending_transactions.id, pending_transactions.type, pending_transactions.currency, pending_transactions.amount, pending_transactions.description, pending_transactions.created_at, actors.id, actors.name FROM pending_transactions, actors WHERE pending_transactions.actor = actors.id ORDER BY pending_transactions.id;")
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		transaction := types.PendingTransaction{}
		err = rows.Scan(&transaction.Id, &transaction.Type, &transaction.Currency, &transaction.Amount, &transaction.Description, &transaction.CreatedAt, &transaction.Actor.Id, &transaction.Actor.Name)
		if err != nil {
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

func CreatePendingTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	transaction := types.PendingTransaction{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data recibida no corresponde con una transacción pendiente")
		return
	}
	if transaction.Type != "input" && transaction.Type != "output" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El tipo de transacción solo puede ser del tipo 'input' o 'output'")
		return
	}
	if transaction.Currency != "USD" && transaction.Currency != "VES" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Solo se aceptan monedas de tipo VES y USD")
		return
	}
	if transaction.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El monto de la transacción pendiente es menor a cero (0)")
		return
	}
	if transaction.Amount > float32(types.MaxTransactionAmount) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El monto de la transacción pendiente exede el máximo permitido")
		return
	}
	if transaction.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente debe poseer una descripción")
		return
	}
	if transaction.Actor.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente debe poseer un actor")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var actorId int
	getActorIdQuery := fmt.Sprintf("SELECT id FROM actors WHERE id='%v';", transaction.Actor.Id)
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

	var insertedId int
	insertTransactionQuery := fmt.Sprintf("INSERT INTO pending_transactions(type, currency, amount, description, actor) VALUES ('%v', '%v', '%v', '%v', '%v') RETURNING id;", transaction.Type, transaction.Currency, transaction.Amount, transaction.Description, transaction.Actor.Id)

	rowsInsertedId, err := db.Query(insertTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	for rowsInsertedId.Next() {
		err = rowsInsertedId.Scan(&insertedId)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	insertedTransaction := types.PendingTransaction{}
	retrieveTransactionQuery := fmt.Sprintf("SELECT pending_transactions.id, pending_transactions.type, pending_transactions.currency, pending_transactions.amount, pending_transactions.description, pending_transactions.created_at, actors.id, actors.name FROM pending_transactions, actors WHERE pending_transactions.actor = actors.id AND pending_transactions.id = '%v';", insertedId)
	rowsRetrievedTransaction, err := db.Query(retrieveTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsRetrievedTransaction.Close()
	for rowsRetrievedTransaction.Next() {
		err = rowsRetrievedTransaction.Scan(&insertedTransaction.Id, &insertedTransaction.Type, &insertedTransaction.Currency, &insertedTransaction.Amount, &insertedTransaction.Description, &insertedTransaction.CreatedAt, &insertedTransaction.Actor.Id, &insertedTransaction.Actor.Name)
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

func PatchPendingTransaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedId := ps.ByName("id")
	pendingTransactionsId, err := strconv.Atoi(requestedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	newPendingTransaction := types.PendingTransaction{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
	}
	err = json.Unmarshal(body, &newPendingTransaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data enviada no corresponde con una transacción pendiente parcial")
		return
	}
	if pendingTransactionsId <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No puede modificar la transacción pendiente cero")
		return
	}
	if newPendingTransaction.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente debe poseer una descripión")
		return
	}
	if newPendingTransaction.Type != "input" && newPendingTransaction.Type != "output" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El tipo de la transacción debe ser 'input' o 'output'")
		return
	}
	if newPendingTransaction.Currency != "USD" && newPendingTransaction.Currency != "VES" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Solo se aceptan monedas de tipo VES y USD")
		return
	}
	if newPendingTransaction.Amount <= 0 || newPendingTransaction.Amount > float32(types.MaxTransactionAmount) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El monto de la transacción es muy bajo o muy alto")
		return
	}
	if newPendingTransaction.Actor.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente debe poseer un actor")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var actorId int
	getActorIdQuery := fmt.Sprintf("SELECT id FROM actors WHERE id = '%v';", newPendingTransaction.Actor.Id)
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

	var updatedId int
	updateQuery := fmt.Sprintf("UPDATE pending_transactions SET type='%v', currency='%v', amount='%v', description='%v', actor='%v' WHERE id='%v' RETURNING id;", newPendingTransaction.Type, newPendingTransaction.Currency, newPendingTransaction.Amount, newPendingTransaction.Description, newPendingTransaction.Actor.Id, pendingTransactionsId)
	rowsUpdatedId, err := db.Query(updateQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsUpdatedId.Close()
	for rowsUpdatedId.Next() {
		err = rowsUpdatedId.Scan(&updatedId)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	if updatedId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente con el id %v no existe", pendingTransactionsId)
		return
	}

	modifiedPendingTransaction := types.PendingTransaction{}
	retrieveTransactionQuery := fmt.Sprintf("SELECT pending_transactions.id, pending_transactions.type, pending_transactions.Currency, pending_transactions.amount, pending_transactions.description, pending_transactions.created_at, actors.id, actors.name FROM pending_transactions, actors WHERE pending_transactions.actor = actors.id AND pending_transactions.id = '%v';", updatedId)
	rowsRetrievedTransaction, err := db.Query(retrieveTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsRetrievedTransaction.Close()
	for rowsRetrievedTransaction.Next() {
		err = rowsRetrievedTransaction.Scan(&modifiedPendingTransaction.Id, &modifiedPendingTransaction.Type, &modifiedPendingTransaction.Currency, &modifiedPendingTransaction.Amount, &modifiedPendingTransaction.Description, &modifiedPendingTransaction.CreatedAt, &modifiedPendingTransaction.Actor.Id, &modifiedPendingTransaction.Actor.Name)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	response, err := json.Marshal(modifiedPendingTransaction)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func DeletePendingTransaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedId := ps.ByName("id")
	if requestedId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el parametro id en la petición de borrado")
		return
	}
	id, err := strconv.Atoi(requestedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	if id <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No puede modificar la transacción pendiente cero")
		return
	}
	db := database.ConnectDB()
	defer db.Close()
	query := fmt.Sprintf("DELETE FROM pending_transactions WHERE id='%v' RETURNING id;", id)
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	deletedId := types.IdResponse{}
	for rows.Next() {
		err = rows.Scan(&deletedId.Id)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if deletedId.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente con el id %v no existe", requestedId)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(deletedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}

func ExecutePendingTransaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedId := ps.ByName("id")
	if requestedId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el parametro id en la petición de borrado")
		return
	}
	id, err := strconv.Atoi(requestedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	if id <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No puede modificar la transacción pendiente cero")
		return
	}
	db := database.ConnectDB()
	defer db.Close()
	query := fmt.Sprintf("SELECT * FROM pending_transactions WHERE id = '%v';", requestedId)
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	pendingTransaction := types.PendingTransaction{}
	for rows.Next() {
		err = rows.Scan(&pendingTransaction.Id, &pendingTransaction.Type, &pendingTransaction.Currency, &pendingTransaction.Amount, &pendingTransaction.Description, &pendingTransaction.Actor.Id, &pendingTransaction.CreatedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	if pendingTransaction.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente con el id %v no existe", requestedId)
		return
	}

	var lastUSDBalance float32
	var lastVESBalance float32
	getLastBalanceQuery := "SELECT USD_balance, VES_balance FROM transactions_with_balances ORDER BY id desc LIMIT 1;"
	lastTransactionRow, err := db.Query(getLastBalanceQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer lastTransactionRow.Close()
	for lastTransactionRow.Next() {
		if err := lastTransactionRow.Scan(&lastUSDBalance, &lastVESBalance); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	newUSDBalance := lastUSDBalance
	newVESBalance := lastVESBalance

	if pendingTransaction.Currency == "USD" {
		if pendingTransaction.Type == "input" {
			newUSDBalance = lastUSDBalance + pendingTransaction.Amount
			if newUSDBalance > float32(types.MaxBalanceAmount) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Su transacción no pudo ser ejecutada porque excede el balance máximo permitido")
				return
			}
		} 

		if pendingTransaction.Type == "output" {
			newUSDBalance = lastUSDBalance - pendingTransaction.Amount
			if newUSDBalance < 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Su transacción pendiente de id %v no pudo ser ejecutada porque genera un balance menor a cero (0)", requestedId)
				return
			}
		}
	}

	if pendingTransaction.Currency == "VES" {
		if pendingTransaction.Type == "input" {
			newVESBalance = lastVESBalance + pendingTransaction.Amount
			if newVESBalance > float32(types.MaxBalanceAmount) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Su transacción no pudo ser ejecutada porque excede el balance máximo permitido")
				return
			}
		} 

		if pendingTransaction.Type == "output" {
			newVESBalance = lastVESBalance - pendingTransaction.Amount
			if newVESBalance < 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Su transacción pendiente de id %v no pudo ser ejecutada porque genera un balance menor a cero (0)", requestedId)
				return
			}
		}
	}

	var insertedTransactionId int
	insertTransactionQuery := fmt.Sprintf("INSERT INTO transactions_with_balances(type, currency, amount, description, USD_balance, VES_balance, actor, created_at) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v') RETURNING id;", pendingTransaction.Type, pendingTransaction.Currency, pendingTransaction.Amount, pendingTransaction.Description, newUSDBalance, newVESBalance, pendingTransaction.Actor.Id, pendingTransaction.CreatedAt)

	rowsId, err := db.Query(insertTransactionQuery)
	if err != nil {
		if err.Error() == `pq: new row for relation "transactions_with_balances" violates check constraint "transactions_with_balances_usd_balance_check"` {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Su transacción pendiente de id %v no pudo ser ejecutada porque genera un balance menor a cero (0)", requestedId)
			return
		}
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsId.Close()

	for rowsId.Next() {
		if err := rowsId.Scan(&insertedTransactionId); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	deleteQuery := fmt.Sprintf("DELETE FROM pending_transactions WHERE id='%v' RETURNING id;", requestedId)
	_, err = db.Query(deleteQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	insertedTransaction := types.TransactionWithBalance{}
	retrieveTransactionQuery := fmt.Sprintf("SELECT transactions_with_balances.id, transactions_with_balances.type, transactions_with_balances.currency, transactions_with_balances.amount, transactions_with_balances.description, transactions_with_balances.USD_balance, transactions_with_balances.VES_balance, transactions_with_balances.executed, transactions_with_balances.created_at, actors.id, actors.name FROM transactions_with_balances, actors WHERE transactions_with_balances.actor = actors.id AND transactions_with_balances.id = '%v';", insertedTransactionId)
	rowsRetrievedTransaction, err := db.Query(retrieveTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsRetrievedTransaction.Close()
	for rowsRetrievedTransaction.Next() {
		err = rowsRetrievedTransaction.Scan(&insertedTransaction.Id, &insertedTransaction.Type, &insertedTransaction.Currency, &insertedTransaction.Amount, &insertedTransaction.Description, &insertedTransaction.USDBalance, &insertedTransaction.VESBalance, &insertedTransaction.Executed, &insertedTransaction.CreatedAt, &insertedTransaction.Actor.Id, &insertedTransaction.Actor.Name)
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

func GetLastTransactionId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	lastPendingTransactionsId := types.IdResponse{
		Id: -1,
	}
	db := database.ConnectDB()
	defer db.Close()
	query := "SELECT id FROM pending_transactions ORDER BY id DESC LIMIT 1;"
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&lastPendingTransactionsId.Id); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if lastPendingTransactionsId.Id == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No existen más transacciones")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(lastPendingTransactionsId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}
