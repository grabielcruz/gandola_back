package main

import (
	"log"
	"net/http"

	"example.com/backend_gandola_soft/actors"
	"example.com/backend_gandola_soft/bills"
	"example.com/backend_gandola_soft/handle_uploads"
	"example.com/backend_gandola_soft/notes"
	"example.com/backend_gandola_soft/pending_transactions"
	"example.com/backend_gandola_soft/transactions"

	"github.com/julienschmidt/httprouter"
)

func CustomOptions(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Enable Cors
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h(w, r, ps)
	}
}

func main() {
	router := httprouter.New()

	router.GET("/", CustomOptions(transactions.Index))

	router.GET("/transactions", CustomOptions(transactions.GetTransactions))
	router.POST("/transactions", CustomOptions(transactions.CreateTransaction))
	router.PATCH("/transactions/:id", CustomOptions(transactions.PatchTransaction))
	router.DELETE("/transactions", CustomOptions(transactions.DeleteLastTransaction))
	router.PUT("/transactions", CustomOptions(transactions.UnexecuteLastTransaction))

	router.GET("/pending_transactions", CustomOptions(pending_transactions.GetPendingTransactions))
	router.POST("/pending_transactions", CustomOptions(pending_transactions.CreatePendingTransaction))
	router.PATCH("/pending_transactions/:id", CustomOptions(pending_transactions.PatchPendingTransaction))
	router.DELETE("/pending_transactions/:id", CustomOptions(pending_transactions.DeletePendingTransaction))
	router.PUT("/pending_transactions/:id", CustomOptions(pending_transactions.ExecutePendingTransaction))

	router.GET("/actors", CustomOptions(actors.GetActors))
	router.POST("/actors", CustomOptions(actors.CreateActor))
	router.PATCH("/actors/:id", CustomOptions(actors.PatchActor))
	router.DELETE("/actors/:id", CustomOptions(actors.DeleteActor))

	router.GET("/notes", CustomOptions(notes.GetNotes))
	router.POST("/notes", CustomOptions(notes.CreateNote))
	router.PATCH("/notes/:id", CustomOptions(notes.PatchNote))
	router.DELETE("/notes/:id", CustomOptions(notes.DeleteNote))
	router.PUT("/attend_note/:id", CustomOptions(notes.AttendNote))
	router.PUT("/unattend_note/:id", CustomOptions(notes.UnattendNote))

	router.GET("/bills", CustomOptions(bills.GetBills))
	router.POST("/bills", CustomOptions(bills.CreateBill))
	router.PATCH("/bills/:id", CustomOptions(bills.PatchBill))
	router.DELETE("/bills/:id", CustomOptions(bills.DeleteBill))//TODO: delete actual image when deleting bill

	router.ServeFiles("/public/*filepath", http.Dir("./public"))
	router.POST("/upload/:id", CustomOptions(handle_uploads.UploadFile))

	log.Fatal(http.ListenAndServe(":8080", router))
}

//TOCONSIDER: maybe I should write tests on demand :D, it takes a hell of time!!!
//TODO: check sql injection protection
//TODO: bills pictures associated with trips, a table for bills, a table for trips, and a table to relationate both of them
