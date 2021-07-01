package main

import (
	"log"
	"net/http"

	"example.com/backend_gandola_soft/actors"
	"example.com/backend_gandola_soft/pending_transactions"
	"example.com/backend_gandola_soft/transactions"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", transactions.Index)

	router.GET("/transactions", transactions.GetTransactions)
	router.POST("/transactions", transactions.CreateTransaction)
	router.PATCH("/transactions", transactions.PatchTransaction)
	router.DELETE("/transactions", transactions.DeleteLastTransaction)
	router.PUT("/transactions", transactions.UnexecuteLastTransaction)
	router.GET("/lasttransactionid", transactions.GetLastTransactionId) //mostly for testing porpuses

	router.GET("/pending_transactions", pending_transactions.GetPendingTransactions)
	router.POST("/pending_transactions", pending_transactions.CreatePendingTransaction)
	router.PATCH("/pending_transactions", pending_transactions.PatchPendingTransaction)
	router.DELETE("/pending_transactions/:id", pending_transactions.DeletePendingTransaction)
	router.PUT("/pending_transactions/:id", pending_transactions.ExecutePendingTransaction)
	router.GET("/lastpendingtransactionid", pending_transactions.GetLastTransactionId) //mostly for testing porpuses

	router.GET("/actors", actors.GetActors)
	router.POST("/actors", actors.CreateActor)

	log.Fatal(http.ListenAndServe(":8080", router))
}

//TOCONSIDER: maybe I should write tests on demand :D, it takes a hell of time!!!
//TODO: make separete file for types; check sql injection protection, evaluate missing tests, for execute pending transaction and unexecute last transaction
//TODO: check amounts with more than 2 decimals
