package transactions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestIndex(t *testing.T) {
	router := httprouter.New()
	router.GET("/", Index)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing status code")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
	}
	t.Log("testing body")
	want := "Server working"
	if string(body) != want {
		t.Errorf("body = %v; want %v", string(body), want)
	}
}

func TestGetTransactions(t *testing.T) {
	router := httprouter.New()
	router.GET("/transactions", GetTransactions)

	req, err := http.NewRequest("GET", "/transactions", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /transactions")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing staus code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing body for transaction zero")
	id := 1
	transaction_type := "zero"
	amount := float32(0)
	description := "transaction zero"
	balance := float32(0)

	transactions := []TransactionWithBalance{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &transactions)
	if err != nil {
		t.Error("Reponse body does not contain an array of type TransactionWithBalances")
	}

	if id != transactions[0].Id {
		t.Errorf("Id = %v, want %v", id, transactions[0].Id)
	}
	if transaction_type != transactions[0].Type {
		t.Errorf("Type = %v, want %v", transaction_type, transactions[0].Type)
	}
	if amount != transactions[0].Amount {
		t.Errorf("Amount = %v, want %v", amount, transactions[0].Amount)
	}
	if description != transactions[0].Description {
		t.Errorf("Description = %v, want %v", description, transactions[0].Description)
	}
	if balance != transactions[0].Balance {
		t.Errorf("Balance = %v, want %v", balance, transactions[0].Balance)
	}
}

func TestCreateTransaction(t *testing.T) {
	router := httprouter.New()
	router.POST("/transactions", CreateTransaction)

	transactionType := "input"
	transactionAmount := float32(3)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v"
  }
	`, transactionType, transactionAmount, transactionDescription)
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /transactions")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)


	t.Log("testing successful status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}


	t.Log("testing create transaction success")
	transactionResponse := TransactionWithBalance{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &transactionResponse)
	if err != nil {
		t.Error("Reponse body does not contain a TransactionWithBalances type")
	}
	if (transactionResponse.Type != "input") {
		t.Errorf("transactionResponse.Type = %v, want %v", transactionResponse.Type, transactionType)
	}
	if (transactionResponse.Amount != 3) {
		t.Errorf("transactionResponse.Amount = %v, want %v", transactionResponse.Amount, transactionAmount)
	}
	if (transactionResponse.Description != "abc") {
		t.Errorf("transactionResponse.Description = %v, want %v", transactionResponse.Description, transactionDescription)
	}
}

func TestCreateTransactionWithoutType(t *testing.T) {
	router := httprouter.New()
	router.POST("/transactions", CreateTransaction)
	transactionType := ""
	transactionAmount := float32(3)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v"
  }
	`, transactionType, transactionAmount, transactionDescription)

	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /transactions")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing error message")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	errMessage := "Debe especificar el tipo de transacción"
	if (string(body) != errMessage) {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreateTransactionWithoutAmount(t *testing.T) {
	router := httprouter.New()
	router.POST("/transactions", CreateTransaction)
	transactionType := "input"
	transactionAmount := float32(0)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v"
  }
	`, transactionType, transactionAmount, transactionDescription)

	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /transactions")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing error message")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	errMessage := "El monto de la transacción debe ser mayor a cero"
	if (string(body) != errMessage) {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreateTransactionWithoutDescription(t *testing.T) {
	router := httprouter.New()
	router.POST("/transactions", CreateTransaction)
	transactionType := "input"
	transactionAmount := float32(3)
	transactionDescription := ""
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v"
  }
	`, transactionType, transactionAmount, transactionDescription)

	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /transactions")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing error message")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	errMessage := "La transacción debe poseer una descripción"
	if (string(body) != errMessage) {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreateTransactionWithBadJson(t *testing.T) {
	router := httprouter.New()
	router.POST("/transactions", CreateTransaction)
	transactionType := "input"
	transactionAmount := float32(3)
	transactionDescription := ""
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
  }
	`, transactionType, transactionAmount, transactionDescription)

	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /transactions")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing error message")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	errMessage := "La data enviada no corresponde con una transacción"
	if (string(body) != errMessage) {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreateTransactionWithBalanceLessThanZero(t *testing.T) {
	router := httprouter.New()
	router.POST("/transactions", CreateTransaction)
	transactionType := "output"
	transactionAmount := float32(999999999999)
	transactionDescription := "balance zero"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v"
  }
	`, transactionType, transactionAmount, transactionDescription)

	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /transactions")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing error message")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	errMessage := "Su transacción no pudo ser ejecutada porque genera un balance menor a cero (0)"
	if (string(body) != errMessage) {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}