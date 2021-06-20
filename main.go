package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/backend_gandola_soft/transactions"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func main() {
	router := httprouter.New()
	router.GET("/", Index)

	router.GET("/transactions", transactions.GetTransactions)
	router.POST("/transactions", transactions.CreateTransaction)
	router.PATCH("/transactions", transactions.PatchTransaction)
	router.DELETE("/transactions", transactions.DeleteLastTransaction) //Delete only the last transaction

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Server working")
}

