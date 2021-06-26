package main

import (
	"log"
	"net/http"

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
	router.GET("/lasttransactionid", transactions.GetLastTransactionId) 

	log.Fatal(http.ListenAndServe(":8080", router))
}

//Todo:use params on patch, explore reset reset serial on DeleteLastTransaction. Proceed with pendings as a regular crud