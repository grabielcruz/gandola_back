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
	router.DELETE("/transactions", transactions.DeleteLastTransaction) //Delete only the last transaction

	log.Fatal(http.ListenAndServe(":8080", router))
}
