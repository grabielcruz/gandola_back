package pending_transactions

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

func TestGetPendingTransactions(t *testing.T) {
	router := httprouter.New()
	router.GET("/pending_transactions", GetPendingTransactions)

	req, err := http.NewRequest("GET", "/pending_transactions", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /pending_transactions")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing for an array of pending transactions")
	pendingTransactions := []PendingTransaction{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &pendingTransactions)
	if err != nil {
		t.Error("Reponse body does not contain an array of type TransactionWithBalances")
	}
}

func TestCreatePendingTransaction(t *testing.T) {
	router := httprouter.New()
	router.POST("/pending_transactions", CreatePendingTransaction)

	transactionType := "input"
	transactionAmount := float32(5)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": 1
  }
	`, transactionType, transactionAmount, transactionDescription)
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /pending_transactions")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing OK status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing create transaction success")
	transactionResponse := PendingTransaction{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &transactionResponse)
	if err != nil {
		t.Error("Response body does not contain a PendingTransaction type")
	}
	if transactionResponse.Type != "input" {
		t.Errorf("transactionResponse.Type = %v, want %v", transactionResponse.Type, transactionType)
	}
	if transactionResponse.Amount != 5 {
		t.Errorf("transactionResponse.Amount = %v, want %v", transactionResponse.Amount, transactionAmount)
	}
	if transactionResponse.Description != "abc" {
		t.Errorf("transactionResponse.Description = %v, want %v", transactionResponse.Description, transactionDescription)
	}
}

func TestCreateTransactionWithoutType(t *testing.T) {
	router := httprouter.New()
	router.POST("/pending_transactions", CreatePendingTransaction)

	transactionType := ""
	transactionAmount := float32(5)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": 1
  }
	`, transactionType, transactionAmount, transactionDescription)
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /pending_transactions")
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
	errMessage := "Debe especificar el tipo de transacción pendiente"
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreatePendingTransactionWithWrongType(t *testing.T) {
	router := httprouter.New()
	router.POST("/pending_transactions", CreatePendingTransaction)

	transactionType := "noType"
	transactionAmount := float32(5)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": 1
  }
	`, transactionType, transactionAmount, transactionDescription)
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /pending_transactions")
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
	errMessage := "Debe especificar el tipo de transacción pendiente"
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreatePendingTransactionWithoutAmount(t *testing.T) {
	router := httprouter.New()
	router.POST("/pending_transactions", CreatePendingTransaction)

	transactionType := "input"
	transactionAmount := float32(0)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": 1
  }
	`, transactionType, transactionAmount, transactionDescription)
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /pending_transactions")
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
	errMessage := "El monto de la transacción pendiente debe ser mayor a cero"
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreatePendingTransactionWithoutDescription(t *testing.T) {
	router := httprouter.New()
	router.POST("/pending_transactions", CreatePendingTransaction)

	transactionType := "input"
	transactionAmount := float32(5)
	transactionDescription := ""
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": 1
  }
	`, transactionType, transactionAmount, transactionDescription)
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /pending_transactions")
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
	errMessage := "La transacción pendiente debe poseer una descripción"
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreatePendingTransactionWithBadJson(t *testing.T) {
	router := httprouter.New()
	router.POST("/pending_transactions", CreatePendingTransaction)

	transactionType := "input"
	transactionAmount := float32(5)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": 1,
  }
	`, transactionType, transactionAmount, transactionDescription)
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /pending_transactions")
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
	errMessage := "La data recibida no corresponde con una transacción pendiente"
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreatePendingTransactionWithNonExistingActor(t *testing.T) {
	router := httprouter.New()
	router.POST("/pending_transactions", CreatePendingTransaction)

	transactionType := "input"
	transactionAmount := float32(5)
	transactionDescription := "abc"
	bodyString := fmt.Sprintf(`
	{
    "Type": "%v",
    "Amount": %v,
    "Description": "%v",
		"Actor": 9999
  }
	`, transactionType, transactionAmount, transactionDescription)
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /pending_transactions")
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

func TestPatchPendingTransaction(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastpendingtransactionid", GetLastTransactionId)

	var lastId IdResponse

	req, err := http.NewRequest("GET", "/lastpendingtransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastpendingtransactionid")
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

	router.PATCH("/pending_transactions", PatchPendingTransaction)
	id := lastId.Id
	transactionType := "output"
	amount := float32(5)
	description := "patch pending transaction test"
	bodyString := fmt.Sprintf(`
		{
			"Id": %v,
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": 1
		}
	`, id, transactionType, amount, description )
	transactionBody := strings.NewReader(bodyString)
	req, err = http.NewRequest("PATCH", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
	}

	router.ServeHTTP(rr, req)
	t.Log("testing OK request status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing patch pending transaction success")
	pendingTransactionResponse := PendingTransaction{}
	body, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &pendingTransactionResponse)
	if err != nil {
		t.Error("Reponse body does not contain a PendingTransaction type")
	}

	if pendingTransactionResponse.Id != id {
		t.Errorf("pendingTransactionResponse.Id = %v, want %v", pendingTransactionResponse.Id, id)
	}

	if pendingTransactionResponse.Type != transactionType {
		t.Errorf("pendingTransactionResponse.Type = %v, want %v", pendingTransactionResponse.Type, transactionType)
	}

	if pendingTransactionResponse.Amount != amount {
		t.Errorf("pendingTransactionResponse.Amount = %v, want %v", pendingTransactionResponse.Amount, amount)
	}

	if pendingTransactionResponse.Description != description {
		t.Errorf("pendingTransactionResponse.Id = %v, want %v", pendingTransactionResponse.Description, description)
	}
}

func TestPatchPendingTransactionEmptyDescription(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastpendingtransactionid", GetLastTransactionId)

	var lastId IdResponse

	req, err := http.NewRequest("GET", "/lastpendingtransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastpendingtransactionid")
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

	router.PATCH("/pending_transactions", PatchPendingTransaction)
	id := lastId.Id
	transactionType := "output"
	amount := float32(5)
	description := ""
	bodyString := fmt.Sprintf(`
		{
			"Id": %v,
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": 1
		}
	`, id, transactionType, amount, description )
	transactionBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
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
	expected := "La transacción pendiente debe poseer una descripión"

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

func TestPatchPendingTransactionBadType(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastpendingtransactionid", GetLastTransactionId)

	var lastId IdResponse

	req, err := http.NewRequest("GET", "/lastpendingtransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastpendingtransactionid")
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

	router.PATCH("/pending_transactions", PatchPendingTransaction)
	id := lastId.Id
	transactionType := "noType"
	amount := float32(5)
	description := "bad type"
	bodyString := fmt.Sprintf(`
		{
			"Id": %v,
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": 1
		}
	`, id, transactionType, amount, description )
	transactionBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
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
	expected := "El tipo de la transacción debe ser 'input' o 'output'"

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

func TestPatchPendingTransactionAmountZeroOrLess(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastpendingtransactionid", GetLastTransactionId)

	var lastId IdResponse

	req, err := http.NewRequest("GET", "/lastpendingtransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastpendingtransactionid")
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

	router.PATCH("/pending_transactions", PatchPendingTransaction)
	id := lastId.Id
	transactionType := "input"
	amount := float32(0)
	description := "amount zero"
	bodyString := fmt.Sprintf(`
		{
			"Id": %v,
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": 1
		}
	`, id, transactionType, amount, description )
	transactionBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
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
	expected := "La transacción pendiente debe poseer un monto mayor que cero"

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

func TestPatchPendingTransactionZero(t *testing.T) {
	router := httprouter.New()

	router.PATCH("/pending_transactions", PatchPendingTransaction)
	id := 1
	transactionType := "output"
	amount := float32(5)
	description := "patching transaction zero"
	bodyString := fmt.Sprintf(`
		{
			"Id": %v,
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": 1
		}
	`, id, transactionType, amount, description )
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("PATCH", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create transaction with empty descripiton fail")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	expected := "No puede modificar la transacción pendiente cero"

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

func TestPatchPendingTransactionBadJson(t *testing.T) {
	router := httprouter.New()

	router.PATCH("/pending_transactions", PatchPendingTransaction)
	id := 2
	transactionType := "output"
	amount := float32(5)
	description := "patching transaction zero"
	bodyString := fmt.Sprintf(`
		{
			"Id": %v,
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": 1,
		}
	`, id, transactionType, amount, description )
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("PATCH", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create transaction with empty descripiton fail")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	expected := "La data enviada no corresponde con una transacción pendiente parcial"

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

func TestPatchPendingTransactionNonExistingId(t *testing.T) {
	router := httprouter.New()

	router.PATCH("/pending_transactions", PatchPendingTransaction)
	id := 9999999
	transactionType := "output"
	amount := float32(5)
	description := "non existing id"
	bodyString := fmt.Sprintf(`
		{
			"Id": %v,
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": 1
		}
	`, id, transactionType, amount, description )
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("PATCH", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create transaction with empty descripiton fail")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	expected := fmt.Sprintf("La transacción pendiente con el id %v no existe", id)

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

func TestPatchPendingTransactionNonExistingActor(t *testing.T) {
	router := httprouter.New()

	router.PATCH("/pending_transactions", PatchPendingTransaction)
	id := 9999999
	transactionType := "output"
	amount := float32(5)
	description := "non existing id"
	bodyString := fmt.Sprintf(`
		{
			"Id": %v,
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": 9999
		}
	`, id, transactionType, amount, description )
	transactionBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("PATCH", "/pending_transactions", transactionBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing patch pending transaction with non existing actor")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	expected := "El actor especificado no existe"

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

func TestDeletePendingTransaction(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastpendingtransactionid", GetLastTransactionId)

	var lastId IdResponse

	req, err := http.NewRequest("GET", "/lastpendingtransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastpendingtransactionid")
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

	router.DELETE("/pending_transactions/:id", DeletePendingTransaction)
	urlRequest := fmt.Sprintf("/pending_transactions/%v", lastId.Id)
	req2, err := http.NewRequest("DELETE", urlRequest, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
	}

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing delete pending transaction success")
	deletedId := IdResponse{}
	body, err = ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &deletedId)
	if err != nil {
		t.Error("Reponse body does not contain an IdResponse type")
	}

	if deletedId.Id != lastId.Id {
		t.Errorf("pendingTransactionResponse.Id = %v, want %v", deletedId.Id, lastId.Id)
	}
}

func TestDeletePendingTransactionWithNoId(t *testing.T) {
	router := httprouter.New()

	router.DELETE("/pending_transactions/:id", DeletePendingTransaction)
	req, err := http.NewRequest("DELETE", "/pending_transactions/", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing not found request status code")
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("status = %v, want %v", status, http.StatusNotFound)
	}

	t.Log("testing delete pending transaction fail")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "404 page not found\n"

	if string(body) != wanted {
		t.Errorf("message = %v, want %v", string(body), wanted)
	}
}

func TestDeletePendingTransactionWithWrongParam(t *testing.T) {
	router := httprouter.New()

	router.DELETE("/pending_transactions/:abc", DeletePendingTransaction)
	req, err := http.NewRequest("DELETE", "/pending_transactions/5", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing delete pending transaction fail")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Debe especificar el parametro id en la petición de borrado"

	if string(body) != wanted {
		t.Errorf("message = %v, want %v", string(body), wanted)
	}
}

func TestDeletePendingTransactionWithNonExistingId(t *testing.T) {
	router := httprouter.New()

	router.DELETE("/pending_transactions/:id", DeletePendingTransaction)
	req, err := http.NewRequest("DELETE", "/pending_transactions/99999", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad request request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing delete pending transaction fail")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "La transacción pendiente con el id 99999 no existe"

	if string(body) != wanted {
		t.Errorf("message = %v, want %v", string(body), wanted)
	}
}

func TestDeletePendingTransactionWithIdEqualOrLessThanOne(t *testing.T) {
	router := httprouter.New()

	router.DELETE("/pending_transactions/:id", DeletePendingTransaction)
	req, err := http.NewRequest("DELETE", "/pending_transactions/1", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /pending_transactions")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad request request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing delete pending transaction fail")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "No puede modificar la transacción pendiente cero"

	if string(body) != wanted {
		t.Errorf("message = %v, want %v", string(body), wanted)
	}
}


// func TestExecutePendingTransaction(t *testing.T) {
// 	router := httprouter.New()
// 	router.GET("/lastpendingtransactionid", GetLastTransactionId)

// 	var lastId IdResponse

// 	req, err := http.NewRequest("GET", "/lastpendingtransactionid", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Error("Could not make a get request to /lastpendingtransactionid")
// 	}
// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	t.Log("testing OK request status code for getting last id")
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("status = %v, want %v", status, http.StatusOK)
// 	}

// 	body, err := ioutil.ReadAll(rr.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Error("Could not read body of response")
// 	}

// 	err = json.Unmarshal(body, &lastId)
// 	if err != nil {
// 		t.Error("Could not read last id from response")
// 	}

// 	router.PUT("/pending_transactions/:id", ExecutePendingTransaction)
// 	requestUrl := fmt.Sprintf("/pending_transactions/%v", lastId.Id)
// 	req2, err := http.NewRequest("PUT", requestUrl, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Errorf("Could not make a patch request to /pending_transactions/%v", lastId.Id)
// 	}

// 	rr2 := httptest.NewRecorder()
// 	router.ServeHTTP(rr2, req2)
// 	t.Log("testing OK request status code")
// 	if status := rr2.Code; status != http.StatusOK {
// 		t.Errorf("status = %v, want %v", status, http.StatusOK)
// 	}

// 	t.Log("testing execute pending transaction success")
// 	insertedTransaction := transactions.TransactionWithBalance{}
// 	body, err = ioutil.ReadAll(rr2.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Error("Could not read body of response")
// 	}
// 	err = json.Unmarshal(body, &insertedTransaction)
// 	if err != nil {
// 		t.Error("Reponse body does not contain a TransactionWithBalance type")
// 	}
// }