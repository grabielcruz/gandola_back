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
	router.PATCH("/transactions/:id", transactions.PatchTransaction)
	router.DELETE("/transactions", transactions.DeleteLastTransaction)
	router.PUT("/transactions", transactions.UnexecuteLastTransaction)
	router.GET("/lasttransactionid", transactions.GetLastTransactionId) //mostly for testing porpuses

	router.GET("/pending_transactions", pending_transactions.GetPendingTransactions)
	router.POST("/pending_transactions", pending_transactions.CreatePendingTransaction)
	router.PATCH("/pending_transactions/:id", pending_transactions.PatchPendingTransaction)
	router.DELETE("/pending_transactions/:id", pending_transactions.DeletePendingTransaction)
	router.PUT("/pending_transactions/:id", pending_transactions.ExecutePendingTransaction)
	router.GET("/lastpendingtransactionid", pending_transactions.GetLastTransactionId) //mostly for testing porpuses

	router.GET("/actors", actors.GetActors)
	router.POST("/actors", actors.CreateActor)
	router.PATCH("/actors/:id", actors.PatchActor)
	router.DELETE("/actors/:id", actors.DeleteActor)
	router.GET("/lastactor", actors.GetLastActor)

	log.Fatal(http.ListenAndServe(":8080", router))
}

//TOCONSIDER: maybe I should write tests on demand :D, it takes a hell of time!!!
//TODO: check sql injection protection
//TODO: bills pictures associated with trips, a table for bills, a table for trips, and a table to relationate both of them