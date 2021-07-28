package bills

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/backend_gandola_soft/types"
	"github.com/julienschmidt/httprouter"
)

func TestGetBills(t *testing.T) {
	router := httprouter.New()
	router.GET("/bills", GetBills)

	req, err := http.NewRequest("GET", "/bills", nil) 
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /bills")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing for an array of bills")
	bills := []types.Bill{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &bills)
	if err != nil {
		t.Error("Response body does not contain an array of type Bill")
	}
}

func TestCreateBill(t *testing.T) {
	router := httprouter.New()
	router.POST("/bills", CreateBill)

	companyId := 2
	date := time.Now().Local().Format(types.DateFormat)
	
	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Date = date
	newBill.Company.Id = companyId
	newBillString, err := json.Marshal(newBill)
	if err != nil {
		t.Log("error on generating new bill")
	}
	requestBody := strings.NewReader(string(newBillString))
	req, err := http.NewRequest("POST", "/bills", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /bills")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing Ok status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing create bill success")
	requestResponse := types.Bill{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &requestResponse)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Bill type")
	}

	if requestResponse.Company.Id != companyId {
		t.Errorf("requestResponse.Company.Id = %v, want %v", requestResponse.Company.Id, companyId)
	}
}

func TestCreateBillWithoutdate(t *testing.T) {
	router := httprouter.New()
	router.POST("/bills", CreateBill)

	companyId := 2
	
	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Company.Id = companyId
	newBillString, err := json.Marshal(newBill)
	if err != nil {
		t.Log("error on generating new bill")
	}
	requestBody := strings.NewReader(string(newBillString))
	req, err := http.NewRequest("POST", "/bills", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /bills")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing Ok status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing create bill success")
	requestResponse := types.Bill{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &requestResponse)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Bill type")
	}

	if requestResponse.Company.Id != companyId {
		t.Errorf("requestResponse.Company.Id = %v, want %v", requestResponse.Company.Id, companyId)
	}
}

func TestCreateBillBaCompanyType(t *testing.T) {
	router := httprouter.New()
	router.POST("/bills", CreateBill)

	companyId := 1
	
	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Company.Id = companyId
	newBillString, err := json.Marshal(newBill)
	if err != nil {
		t.Log("error on generating new bill")
	}
	requestBody := strings.NewReader(string(newBillString))
	req, err := http.NewRequest("POST", "/bills", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /bills")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing Ok status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create bill success")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	wanted := "La compañía especificada no es mina o contratante"
	response := string(body)
	if response != wanted {
		t.Errorf("response = '%v', wanted = '%v'", response, wanted)
	}
}

func TestCreateBillWithoutCompany(t *testing.T) {
	router := httprouter.New()
	router.POST("/bills", CreateBill)

	newBill := types.Bill{}
	newBill.Code = "1234"
	newBillString, err := json.Marshal(newBill)
	if err != nil {
		t.Log("error on generating new bill")
	}
	requestBody := strings.NewReader(string(newBillString))
	req, err := http.NewRequest("POST", "/bills", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /bills")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create bill success")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	response := string(body)
	wanted := "Debe especificar la compañía a la que pertenece la factura"
	if response != wanted {
		t.Errorf("response = %v, want %v", response, wanted)
	}
}

func TestCreateBillNonExistingCompany(t *testing.T) {
	router := httprouter.New()
	router.POST("/bills", CreateBill)

	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Company.Id = 99999
	newBillString, err := json.Marshal(newBill)
	if err != nil {
		t.Log("error on generating new bill")
	}
	requestBody := strings.NewReader(string(newBillString))
	req, err := http.NewRequest("POST", "/bills", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /bills")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create bill success")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	response := string(body)
	wanted := "La compañía especificada no existe"
	if response != wanted {
		t.Errorf("response = %v, want %v", response, wanted)
	}
}

func TestCreateBillBadDate(t *testing.T) {
	router := httprouter.New()
	router.POST("/bills", CreateBill)

	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Date = "bad date"
	newBill.Company.Id = 2
	newBillString, err := json.Marshal(newBill)
	if err != nil {
		t.Log("error on generating new bill")
	}
	requestBody := strings.NewReader(string(newBillString))
	req, err := http.NewRequest("POST", "/bills", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /bills")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create bill success")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	response := string(body)
	wanted := "La fecha de la factura no tiene un formato válido"
	if response != wanted {
		t.Errorf("response = %v, want %v", response, wanted)
	}
}

func TestPatchBill(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastbillid", GetLastBillId)

	var lastBillId types.IdResponse
	req, err := http.NewRequest("GET", "/lastbillid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastbillid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las bill id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastBillId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PATCH("/bills/:id", PatchBill)
	requestUrl := fmt.Sprintf("/bills/%v", lastBillId.Id)
	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Url = "new_photo.jpg"
	newBill.Date = time.Now().Local().Format(types.DateFormat)
	newBill.Company.Id = 2
	newBill.Charged = true
	newBillJson, err := json.Marshal(newBill)
	if err != nil {
		t.Error("Could not marshal new bill into json")
	}
	requestBody := strings.NewReader(string(newBillJson))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /bills/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for patching last note")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing getting back a bill from patch request")
	responseBill := types.Bill{}
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body2, &responseBill)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Bill type")
	}
}

func TestPatchBillBadDate(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastbillid", GetLastBillId)

	var lastBillId types.IdResponse
	req, err := http.NewRequest("GET", "/lastbillid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastbillid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las bill id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastBillId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PATCH("/bills/:id", PatchBill)
	requestUrl := fmt.Sprintf("/bills/%v", lastBillId.Id)
	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Url = "new_photo.jpg"
	newBill.Date = "bad date"
	newBill.Company.Id = 2
	newBill.Charged = true
	newBillJson, err := json.Marshal(newBill)
	if err != nil {
		t.Error("Could not marshal new bill into json")
	}
	requestBody := strings.NewReader(string(newBillJson))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /bills/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a bill from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	
	response := string(body2)
	wanted := "La fecha de la factura no tiene un formato válido"
	if response != wanted {
		t.Errorf("response = '%v', wanted = '%v'", response, wanted)
	}
}

func TestPatchBillNoCompany(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastbillid", GetLastBillId)

	var lastBillId types.IdResponse
	req, err := http.NewRequest("GET", "/lastbillid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastbillid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las bill id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastBillId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PATCH("/bills/:id", PatchBill)
	requestUrl := fmt.Sprintf("/bills/%v", lastBillId.Id)
	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Url = "new_photo.jpg"
	newBill.Date = time.Now().Local().Format(types.DateFormat)
	newBill.Charged = true
	newBillJson, err := json.Marshal(newBill)
	if err != nil {
		t.Error("Could not marshal new bill into json")
	}
	requestBody := strings.NewReader(string(newBillJson))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /bills/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing bad request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a bill from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	response := string(body2)
	wanted := "Debe especificar la compañía a la que pertenece la factura"
	if response != wanted {
		t.Errorf("response = '%v', wanted='%v'", response, wanted)
	}
}

func TestPatchBillBadCompanyType(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastbillid", GetLastBillId)

	var lastBillId types.IdResponse
	req, err := http.NewRequest("GET", "/lastbillid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastbillid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las bill id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastBillId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PATCH("/bills/:id", PatchBill)
	requestUrl := fmt.Sprintf("/bills/%v", lastBillId.Id)
	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Url = "new_photo.jpg"
	newBill.Company.Id = 1
	newBill.Date = time.Now().Local().Format(types.DateFormat)
	newBill.Charged = true
	newBillJson, err := json.Marshal(newBill)
	if err != nil {
		t.Error("Could not marshal new bill into json")
	}
	requestBody := strings.NewReader(string(newBillJson))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /bills/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing bad request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a bill from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	response := string(body2)
	wanted := "La compañía especificada no es mina o contratante"
	if response != wanted {
		t.Errorf("response = '%v', wanted='%v'", response, wanted)
	}
}


func TestPatchBillNonExistingCompany(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastbillid", GetLastBillId)

	var lastBillId types.IdResponse
	req, err := http.NewRequest("GET", "/lastbillid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastbillid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las bill id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastBillId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PATCH("/bills/:id", PatchBill)
	requestUrl := fmt.Sprintf("/bills/%v", lastBillId.Id)
	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Url = "new_photo.jpg"
	newBill.Date = time.Now().Local().Format(types.DateFormat)
	newBill.Company.Id = 9999
	newBill.Charged = true
	newBillJson, err := json.Marshal(newBill)
	if err != nil {
		t.Error("Could not marshal new bill into json")
	}
	requestBody := strings.NewReader(string(newBillJson))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /bills/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing bad request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a bill from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	response := string(body2)
	wanted := "La compañía especificada no existe"
	if response != wanted {
		t.Errorf("response = '%v', wanted='%v'", response, wanted)
	}
}

func TestPatchBillBadId(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastbillid", GetLastBillId)

	router.PATCH("/bills/:id", PatchBill)
	requestUrl := fmt.Sprintf("/bills/%v", 0)
	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Url = "new_photo.jpg"
	newBill.Date = time.Now().Local().Format(types.DateFormat)
	newBill.Charged = true
	newBillJson, err := json.Marshal(newBill)
	if err != nil {
		t.Error("Could not marshal new bill into json")
	}
	requestBody := strings.NewReader(string(newBillJson))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /bills/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing bad request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a bill from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	response := string(body2)
	wanted := "El id de la factura debe ser mayor a cero"
	if response != wanted {
		t.Errorf("response = '%v', wanted='%v'", response, wanted)
	}
}

func TestPatchBillNonExistingId(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastbillid", GetLastBillId)

	router.PATCH("/bills/:id", PatchBill)
	requestUrl := fmt.Sprintf("/bills/%v", 999999)
	newBill := types.Bill{}
	newBill.Code = "1234"
	newBill.Url = "new_photo.jpg"
	newBill.Date = time.Now().Local().Format(types.DateFormat)
	newBill.Company.Id = 2
	newBill.Charged = true
	newBillJson, err := json.Marshal(newBill)
	if err != nil {
		t.Error("Could not marshal new bill into json")
	}
	requestBody := strings.NewReader(string(newBillJson))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /bills/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing bad request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a bill from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	response := string(body2)
	wanted := "La factura solicitada no existe"
	if response != wanted {
		t.Errorf("response = '%v', wanted='%v'", response, wanted)
	}
}

func TestDeleteBill(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastbillid", GetLastBillId)

	var lastBillId types.IdResponse
	req, err := http.NewRequest("GET", "/lastbillid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastbillid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las bill id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastBillId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.DELETE("/bills/:id", DeleteBill)
	requestUrl := fmt.Sprintf("/bills/%v", lastBillId.Id)

	req2, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /bills/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for deleting last note")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing getting back an IdResponse from delete request")
	response := types.IdResponse{}
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body2, &response)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an IdResponse type")
	}

	t.Log("testing sent and received Id are identicals")
	if response.Id != lastBillId.Id {
		t.Errorf("response.Id = '%v', lastBillId.Id = '%v', they are different", response.Id, lastBillId.Id)
	}
}

func TestDeleteBillZeroId(t *testing.T) {
	router := httprouter.New()

	router.DELETE("/bills/:id", DeleteBill)
	requestUrl := fmt.Sprintf("/bills/%v", 0)

	req2, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /bills/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for deleting last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an IdResponse from delete request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	response := string(body2)
	wanted := "El id de la factura debe ser mayor a cero"
	if response != wanted {
		t.Errorf("response = '%v', want = '%v'", response, wanted)
	}
}

func TestDeleteBillBadId(t *testing.T) {
	router := httprouter.New()

	router.DELETE("/bills/:id", DeleteBill)
	requestUrl := fmt.Sprintf("/bills/%v", 9999)

	req2, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /bills/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for deleting last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an IdResponse from delete request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	response := string(body2)
	wanted := "La factura con el id 9999 no existe"
	if response != wanted {
		t.Errorf("response = '%v', want = '%v'", response, wanted)
	}
}