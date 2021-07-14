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

	"example.com/backend_gandola_soft/types"
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
	pendingTransactions := []types.PendingTransaction{}
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
		"Actor": {
			"Id": 1
		}
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
	transactionResponse := types.PendingTransaction{}
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
		"Actor": {
			"Id": 1
		}
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
		"Actor": {
			"Id": 1
		}
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
		"Actor": {
			"Id": 1
		}
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
	errMessage := "El monto de la transacción pendiente es menor a cero (0)"
	if string(body) != errMessage {
		t.Errorf("response = %v, want %v", string(body), errMessage)
	}
}

func TestCreatePendingTransactionMoreThanMaximum(t *testing.T) {
	router := httprouter.New()
	router.POST("/pending_transactions", CreatePendingTransaction)

	transactionType := "input"
	transactionAmount := float32(1e15)
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
	errMessage := "El monto de la transacción pendiente exede el máximo permitido"
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
		"Actor": {
			"Id": 1
		}
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
		"Actor": {
			"Id": 1
		},
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
		"Actor": {
			"Id": 9999
		}
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

	var lastId types.IdResponse

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

	router.PATCH("/pending_transactions/:id", PatchPendingTransaction)
	id := lastId.Id
	transactionType := "output"
	amount := float32(5)
	description := "patch pending transaction test"
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": {
				"Id": 1
			}
		}
	`, transactionType, amount, description)
	transactionBody := strings.NewReader(bodyString)
	urlRequest := fmt.Sprintf("/pending_transactions/%v", id)
	req, err = http.NewRequest("PATCH", urlRequest, transactionBody)
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
	pendingTransactionResponse := types.PendingTransaction{}
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

	var lastId types.IdResponse

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

	router.PATCH("/pending_transactions/:id", PatchPendingTransaction)
	id := lastId.Id
	transactionType := "output"
	amount := float32(5)
	description := ""
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": {
				"Id": 1
			}
		}
	`, transactionType, amount, description)
	transactionBody := strings.NewReader(bodyString)
	urlRequest := fmt.Sprintf("/pending_transactions/%v", id)
	req2, err := http.NewRequest("PATCH", urlRequest, transactionBody)
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

	var lastId types.IdResponse

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

	router.PATCH("/pending_transactions/:id", PatchPendingTransaction)
	id := lastId.Id
	transactionType := "noType"
	amount := float32(5)
	description := "bad type"
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": {
				"Id": 1
			}
		}
	`, transactionType, amount, description)
	transactionBody := strings.NewReader(bodyString)
	urlRequest := fmt.Sprintf("/pending_transactions/%v", id)
	req2, err := http.NewRequest("PATCH", urlRequest, transactionBody)
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

	var lastId types.IdResponse

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

	router.PATCH("/pending_transactions/:id", PatchPendingTransaction)
	id := lastId.Id
	transactionType := "input"
	amount := float32(0)
	description := "amount zero"
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": {
				"Id": 1
			}
		}
	`, transactionType, amount, description)
	transactionBody := strings.NewReader(bodyString)
	urlRequest := fmt.Sprintf("/pending_transactions/%v", id)
	req2, err := http.NewRequest("PATCH", urlRequest, transactionBody)
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
	expected := "El monto de la transacción es muy bajo o muy alto"

	if string(body) != expected {
		t.Errorf("body = %v, want %v", string(body), expected)
	}
}

func TestPatchPendingTransactionZero(t *testing.T) {
	router := httprouter.New()

	router.PATCH("/pending_transactions/:id", PatchPendingTransaction)
	id := 1
	transactionType := "output"
	amount := float32(5)
	description := "patching transaction zero"
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": {
				"Id": 1
			}
		}
	`, transactionType, amount, description)
	transactionBody := strings.NewReader(bodyString)
	urlRequest := fmt.Sprintf("/pending_transactions/%v", id)
	req, err := http.NewRequest("PATCH", urlRequest, transactionBody)
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

	router.PATCH("/pending_transactions/:id", PatchPendingTransaction)
	id := 2
	transactionType := "output"
	amount := float32(5)
	description := "patching transaction zero"
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": {
				"Id": 1
			},
		}
	`, transactionType, amount, description)
	transactionBody := strings.NewReader(bodyString)
	urlRequest := fmt.Sprintf("/pending_transactions/%v", id)
	req, err := http.NewRequest("PATCH", urlRequest, transactionBody)
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

	router.PATCH("/pending_transactions/:id", PatchPendingTransaction)
	id := 9999999
	transactionType := "output"
	amount := float32(5)
	description := "non existing id"
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": {
				"Id": 1
			}
		}
	`, transactionType, amount, description)
	transactionBody := strings.NewReader(bodyString)
	urlRequest := fmt.Sprintf("/pending_transactions/%v", id)
	req, err := http.NewRequest("PATCH", urlRequest, transactionBody)
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

	router.PATCH("/pending_transactions/:id", PatchPendingTransaction)
	id := 9999999
	transactionType := "output"
	amount := float32(5)
	description := "non existing id"
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Amount": %v,
			"Description": "%v",
			"Actor": {
				"Id": 9999
			}
		}
	`, transactionType, amount, description)
	transactionBody := strings.NewReader(bodyString)
	urlRequest := fmt.Sprintf("/pending_transactions/%v", id)
	req, err := http.NewRequest("PATCH", urlRequest, transactionBody)
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

	var lastId types.IdResponse

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
	deletedId := types.IdResponse{}
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

func TestExecutePendingTransaction(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastpendingtransactionid", GetLastTransactionId)
	router.POST("/pending_transactions", CreatePendingTransaction)

	transactionType := "input"
	transactionAmount := float32(42)
	transactionDescription := "pending transaction for make execution"
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
	transactionResponse := types.PendingTransaction{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &transactionResponse)
	if err != nil {
		t.Error("Response body does not contain a PendingTransaction type")
	}

	router.PUT("/pending_transactions/:id", ExecutePendingTransaction)
	id := transactionResponse.Id

	urlRequest := fmt.Sprintf("/pending_transactions/%v", id)
	req2, err := http.NewRequest("PUT", urlRequest, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a put request to /pending_transactions/:id")
	}

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing put pending transaction success")
	newTransaction := types.TransactionWithBalance{}
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body2, &newTransaction)
	if err != nil {
		t.Error("Reponse body does not contain a TransactionWithBalance type")
	}

	var lastIdAfterExecution types.IdResponse
	req3, err := http.NewRequest("GET", "/lastpendingtransactionid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastpendingtransactionid")
	}
	rr3 := httptest.NewRecorder()
	router.ServeHTTP(rr3, req3)

	t.Log("testing OK request status code for getting last id after execution")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body3, err := ioutil.ReadAll(rr3.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing IdResponse after execution")
	err = json.Unmarshal(body3, &lastIdAfterExecution)
	if err != nil {
		t.Error("Could not read last id from response")
	}

	t.Log("testing IdResponse.Id after execution is less than IdResponse.Id before execution")
	if lastIdAfterExecution.Id >= transactionResponse.Id {
		t.Errorf("lastIdAfterExecution = %v, lastIdBeforeExecution = %v", lastIdAfterExecution.Id, transactionResponse.Id)
	}
}
