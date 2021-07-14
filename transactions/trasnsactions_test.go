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

	"example.com/backend_gandola_soft/types"
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
	transaction_type := "input"
	amount := float32(0)
	description := "transaction zero"
	balance := float32(0)

	transactions := []types.TransactionWithBalance{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &transactions)
	if err != nil {
		t.Error("Reponse body does not contain an array of type TransactionWithBalances")
	}

	lastIndex := len(transactions) - 1
	if id != transactions[lastIndex].Id {
		t.Errorf("Id = %v, want %v", id, transactions[lastIndex].Id)
	}
	if transaction_type != transactions[lastIndex].Type {
		t.Errorf("Type = %v, want %v", transaction_type, transactions[lastIndex].Type)
	}
	if amount != transactions[lastIndex].Amount {
		t.Errorf("Amount = %v, want %v", amount, transactions[lastIndex].Amount)
	}
	if description != transactions[lastIndex].Description {
		t.Errorf("Description = %v, want %v", description, transactions[lastIndex].Description)
	}
	if balance != transactions[lastIndex].Balance {
		t.Errorf("Balance = %v, want %v", balance, transactions[lastIndex].Balance)
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
    "Description": "%v",
		"Actor": {
			"Id": 1
		}
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
	transactionResponse := types.TransactionWithBalance{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &transactionResponse)
	if err != nil {
		t.Error("Reponse body does not contain a TransactionWithBalances type")
	}
	if transactionResponse.Type != "input" {
		t.Errorf("transactionResponse.Type = %v, want %v", transactionResponse.Type, transactionType)
	}
	if transactionResponse.Amount != 3 {
		t.Errorf("transactionResponse.Amount = %v, want %v", transactionResponse.Amount, transactionAmount)
	}
	if transactionResponse.Description != "abc" {
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
    "Description": "%v",
		"Actor": {
			"Id": 1
		}
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
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreateTransactionWithWrongType(t *testing.T) {
	router := httprouter.New()
	router.POST("/transactions", CreateTransaction)
	transactionType := "wrongtype"
	transactionAmount := float32(3)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": {
			"Id": 1
		}
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
	errMessage := "El tipo de transacción solo puede ser del tipo 'input' o 'output'"
	if string(body) != errMessage {
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
    "Description": "%v",
		"Actor": {
			"Id": 1
		}
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
	errMessage := "El monto de la transacción es menor a cero (0)"
	if string(body) != errMessage {
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
    "Description": "%v",
		"Actor": {
			"Id": 1
		}
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
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreateTransactionWithBadJson(t *testing.T) {
	router := httprouter.New()
	router.POST("/transactions", CreateTransaction)
	transactionType := "input"
	transactionAmount := float32(3)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": {
			"Id": 1
		},
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
	errMessage := "La data recibida no corresponde con una transacción"
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreateTransactionWithNonExistingActor(t *testing.T) {
	router := httprouter.New()
	router.POST("/transactions", CreateTransaction)
	transactionType := "input"
	transactionAmount := float32(3)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": {
			"Id": 9999
		}
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
	errMessage := "El actor especificado no existe"
	if string(body) != errMessage {
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
    "Description": "%v",
		"Actor": {
			"Id": 1
		}
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
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreateTransactionMoreThanMaximum(t *testing.T) {
	router := httprouter.New()
	router.POST("/transactions", CreateTransaction)
	transactionType := "input"
	transactionAmount := float32(1e15)
	transactionDescription := "balance zero"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": {
			"Id": 1
		}
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
	errMessage := "El monto de la transacción exede el máximo permitido"
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestPatchTransaction(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttransactionid", GetLastTransactionId)

	var lastId types.IdResponse

	req, err := http.NewRequest("GET", "/lasttransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lasttransactionid")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &lastId)
	if err != nil {
		t.Error("Could not read last id from response")
	}
	router.PATCH("/transactions/:id", PatchTransaction)

	id := lastId.Id
	description := "transaction patch testing"
	bodyString := fmt.Sprintf(`
		{
			"Description": "%v"
		}
	`, description)
	transactionBody := strings.NewReader(bodyString)
	requestUrl := fmt.Sprintf("/transactions/%v", id)
	req2, err := http.NewRequest("PATCH", requestUrl, transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /transactions")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing patch transaction success")
	transactionResponse := types.TransactionWithBalance{}
	body, err = ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &transactionResponse)
	if err != nil {
		t.Error("Reponse body does not contain a TransactionWithBalances type")
	}

	if transactionResponse.Id != id {
		t.Errorf("transactionResponse.Id = %v, want %v", transactionResponse.Id, id)
	}

	if transactionResponse.Description != description {
		t.Errorf("transactionResponse.Id = %v, want %v", transactionResponse.Description, description)
	}
}

func TestPatchTransactionEmptyDescription(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttransactionid", GetLastTransactionId)

	var lastId types.IdResponse

	req, err := http.NewRequest("GET", "/lasttransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lasttransactionid")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &lastId)
	if err != nil {
		t.Error("Could not read last id from response")
	}
	router.PATCH("/transactions/:id", PatchTransaction)

	id := lastId.Id

	bodyString := `
		{
			"Description": ""
		}
	`
	transactionBody := strings.NewReader(bodyString)
	requestUrl := fmt.Sprintf("/transactions/%v", id)
	req2, err := http.NewRequest("PATCH", requestUrl, transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /transactions")
	}

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing bad request status code")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create transaction with empty descripiton fail")
	body, err = ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	expected := "La transacción debe poseer una descripión"

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

func TestPatchTransactionZero(t *testing.T) {
	router := httprouter.New()
	router.PATCH("/transactions/:id", PatchTransaction)

	id := 1
	description := "transaction patch testing"
	bodyString := fmt.Sprintf(`
		{
			"Description": "%v"
		}
	`, description)
	transactionBody := strings.NewReader(bodyString)
	requestUrl := fmt.Sprintf("/transactions/%v", id)
	req, err := http.NewRequest("PATCH", requestUrl, transactionBody)
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

	t.Log("testing attempt to modify transaction zero")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	expected := "No puede modificar la transacción cero"

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

func TestPatchTransactionBadJson(t *testing.T) {
	router := httprouter.New()
	router.PATCH("/transactions/:id", PatchTransaction)

	id := 1
	description := "transaction patch testing"
	bodyString := fmt.Sprintf(`
		{
			"Description": "%v",
		}
	`, description)
	transactionBody := strings.NewReader(bodyString)
	requestUrl := fmt.Sprintf("/transactions/%v", id)
	req, err := http.NewRequest("PATCH", requestUrl, transactionBody)
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

	t.Log("testing attempt to modify transaction zero")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	expected := "La data enviada no corresponde con una transacción parcial"

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)

	}
}

func TestPatchTransactionNonExistingId(t *testing.T) {
	router := httprouter.New()
	router.PATCH("/transactions/:id", PatchTransaction)

	id := 99999
	description := "transaction patch testing"
	bodyString := fmt.Sprintf(`
		{
			"Description": "%v"
		}
	`, description)
	transactionBody := strings.NewReader(bodyString)
	requestUrl := fmt.Sprintf("/transactions/%v", id)
	req, err := http.NewRequest("PATCH", requestUrl, transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /transactions")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing attempt to modify an unexisting transaction")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	expected := fmt.Sprintf("La transacción con el id %v no existe", id)

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

var lastIdBeforeDeletion types.IdResponse

func TestDeleteLastTransaction(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttransactionid", GetLastTransactionId)

	var lastId types.IdResponse
	var deletedId types.IdResponse

	req, err := http.NewRequest("GET", "/lasttransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lasttransactionid")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &lastId)
	if err != nil {
		t.Error("Could not read last id from response")
	}

	router.DELETE("/transactions", DeleteLastTransaction)
	req, err = http.NewRequest("DELETE", "/transactions", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /transactions")
	}
	router.ServeHTTP(rr, req)
	t.Log("testing OK request status code for deleting last id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}
	body, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &deletedId)
	if err != nil {
		t.Error("Could not read deleted id from response")
	}
	lastIdBeforeDeletion.Id = deletedId.Id
	t.Log("Testing last id equals deleted id")
	if lastId.Id != deletedId.Id {
		t.Errorf("deleted id = %v, want %v", deletedId.Id, lastId.Id)
	}
}

func TestRollbackIdOnDelete(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttransactionid", GetLastTransactionId)

	var lastId types.IdResponse

	req, err := http.NewRequest("GET", "/lasttransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lasttransactionid")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &lastId)
	if err != nil {
		t.Error("Could not read last id from response")
	}
	wantedId := lastIdBeforeDeletion.Id - 1 // because we deleted one transaction

	t.Log("testing if last id rolledb back on deletion")
	if lastId.Id != wantedId {
		t.Errorf("last id = %v, want %v", lastId.Id, wantedId)
	}
}

func TestUnexecuteLastTransaction(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttransactionid", GetLastTransactionId)
	router.POST("/transactions", CreateTransaction)

	transactionType := "input"
	transactionAmount := float32(42)
	transactionDescription := "transaction to test unexecution"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": {
			"Id": 1
		}
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
	transactionResponse := types.TransactionWithBalance{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &transactionResponse)
	if err != nil {
		t.Error("Reponse body does not contain a TransactionWithBalances type")
	}

	router.PUT("/transactions", UnexecuteLastTransaction)
	req2, err := http.NewRequest("PUT", "/transactions", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a put request to /transactions")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("Testing Ok status request for unexecuting last transaction")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("Testing pending transaction generated")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of second response")
	}
	newPendingTransaction := types.PendingTransaction{}
	err = json.Unmarshal(body2, &newPendingTransaction)
	if err != nil {
		t.Error("Could not read a pending transaction from response")
	}

	var lastIdAfterUnexecution types.IdResponse

	req3, err := http.NewRequest("GET", "/lasttransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lasttransactionid")
	}
	rr3 := httptest.NewRecorder()
	router.ServeHTTP(rr3, req3)

	t.Log("testing OK request status code for getting last id after unexecution")
	if status := rr3.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing IdResponse after unexecuting last transaction")
	body3, err := ioutil.ReadAll(rr3.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body3, &lastIdAfterUnexecution)
	if err != nil {
		t.Error("Could not read last id from response")
	}

	t.Log("testing second IdResponse.Id is less than last IdResponse by one")
	if (transactionResponse.Id - 1 != lastIdAfterUnexecution.Id) {
		t.Errorf("lastIdBeforeUnexecution.Id = %v, lastIdAfterUnexecution.Id = %v", transactionResponse.Id, lastIdAfterUnexecution.Id)
	}
}