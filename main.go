package main

import (
	"log"
	"net/http"

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
	router.GET("/lasttransactionid", transactions.GetLastTransactionId) //mostly for testing porpuses

	router.GET("/pending_transactions", pending_transactions.GetPendingTransactions)

	log.Fatal(http.ListenAndServe(":8080", router))
}