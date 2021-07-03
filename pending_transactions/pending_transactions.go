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
	rows, err := db.Query("SELECT * FROM pending_transactions ORDER BY id ASC;")
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		transaction := types.PendingTransaction{}
		err = rows.Scan(&transaction.Id, &transaction.Type, &transaction.Amount, &transaction.Description, &transaction.Actor, &transaction.CreatedAt)
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
	transaction := types.PartialPendingTransaction{}
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
	if transaction.Actor <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente debe poseer un actor")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var actorId int
	getActorIdQuery := fmt.Sprintf("SELECT id FROM actors WHERE id=%v", transaction.Actor)
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

	insertedTransaction := types.PendingTransaction{}
	insertTransactionQuery := fmt.Sprintf("INSERT INTO pending_transactions(type, amount, description, actor) VALUES ('%v', '%v', '%v', '%v') RETURNING id, type, amount, description, actor, created_at;", transaction.Type, transaction.Amount, transaction.Description, transaction.Actor)

	rows, err := db.Query(insertTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	for rows.Next() {
		err = rows.Scan(&insertedTransaction.Id, &insertedTransaction.Type, &insertedTransaction.Amount, &insertedTransaction.Description, &insertedTransaction.Actor, &insertedTransaction.CreatedAt)
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
	if newPendingTransaction.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente debe poseer un monto mayor que cero")
		return
	}
	if newPendingTransaction.Actor <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente debe poseer un actor")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var actorId int
	getActorIdQuery := fmt.Sprintf("SELECT id FROM actors WHERE id=%v", newPendingTransaction.Actor)
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

	query := fmt.Sprintf("UPDATE pending_transactions SET type='%v', amount='%v', description='%v', actor='%v' WHERE id='%v' RETURNING id, type, amount, description, actor, created_at;", newPendingTransaction.Type, newPendingTransaction.Amount, newPendingTransaction.Description, newPendingTransaction.Actor, pendingTransactionsId)
	modifiedPendingTransaction := types.PendingTransaction{}
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&modifiedPendingTransaction.Id, &modifiedPendingTransaction.Type, &modifiedPendingTransaction.Amount, &modifiedPendingTransaction.Description, &modifiedPendingTransaction.Actor, &modifiedPendingTransaction.CreatedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	if modifiedPendingTransaction.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La transacción pendiente con el id %v no existe", pendingTransactionsId)
		return
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
	query := fmt.Sprintf("SELECT * FROM pending_transactions WHERE id='%v';", requestedId)
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	pendingTransaction := types.PendingTransaction{}
	for rows.Next() {
		err = rows.Scan(&pendingTransaction.Id, &pendingTransaction.Type, &pendingTransaction.Amount, &pendingTransaction.Description, &pendingTransaction.Actor, &pendingTransaction.CreatedAt)
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

	if pendingTransaction.Type == "input" {
		newBalance = lastBalance + pendingTransaction.Amount
	} else if pendingTransaction.Type == "output" {
		newBalance = lastBalance - pendingTransaction.Amount
		if newBalance < 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Su transacción pendiente de id %v no pudo ser ejecutada porque genera un balance menor a cero (0)", requestedId)
			return
		}
	}

	insertedTransaction := types.TransactionWithBalance{}
	insertTransactionQuery := fmt.Sprintf("INSERT INTO transactions_with_balances(type, amount, description, balance, actor, created_at) VALUES ('%v', '%v', '%v', '%v', '%v', '%v') RETURNING id, type, amount, description, balance, actor, executed, created_at;", pendingTransaction.Type, pendingTransaction.Amount, pendingTransaction.Description, newBalance, pendingTransaction.Actor, pendingTransaction.CreatedAt)

	rows, err = db.Query(insertTransactionQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&insertedTransaction.Id, &insertedTransaction.Type, &insertedTransaction.Amount, &insertedTransaction.Description, &insertedTransaction.Balance, &insertedTransaction.Actor, &insertedTransaction.Executed, &insertedTransaction.CreatedAt); err != nil {
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
	query := "SELECT id FROM pending_transactions ORDER BY id desc LIMIT 1;"
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
