package trucks

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
	"example.com/backend_gandola_soft/utils"
	"github.com/julienschmidt/httprouter"
)

func TestGetTrucks(t *testing.T) {
	router := httprouter.New()
	router.GET("/trucks", GetTrucks)

	req, err := http.NewRequest("GET", "/trucks", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /trucks")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing for an array of trucks")
	trucks := []types.Truck{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &trucks)
	if err != nil {
		t.Error("Response body does not contain an array of type Truck")
	}
}

func TestCreateTruck(t *testing.T) {
	router := httprouter.New()
	router.POST("/trucks", CreateTruck)

	newTruck := types.Truck{}
	newTruck.Name = utils.RandStringBytes(10)
	newTruck.Data = "data for the test truck"

	jsonNewTruck, err := json.Marshal(newTruck)
	if err != nil {
		t.Error("Could not marshal json")
	}
	requestBody := strings.NewReader(string(jsonNewTruck))
	req, err := http.NewRequest("POST", "/trucks", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /trucks")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing OK status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing create truck success")
	requestResponse := types.Truck{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &requestResponse)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Truck type")
	}
}

func TestCreateTruckWithoutName(t *testing.T) {
	router := httprouter.New()
	router.POST("/trucks", CreateTruck)

	newTruck := types.Truck{}
	newTruck.Data = "data for the test truck"

	jsonNewTruck, err := json.Marshal(newTruck)
	if err != nil {
		t.Error("Could not marshal json")
	}
	requestBody := strings.NewReader(string(jsonNewTruck))
	req, err := http.NewRequest("POST", "/trucks", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /trucks")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create truck success")
	
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	response := string(body)
	wanted := "Debe especificar el nombre del camión"
	if response != wanted {
		t.Errorf("response = %v, wanted %v", response, wanted)
	}
}

func TestCreateTruckWithoutData(t *testing.T) {
	router := httprouter.New()
	router.POST("/trucks", CreateTruck)

	newTruck := types.Truck{}
	newTruck.Name = utils.RandStringBytes(10)

	jsonNewTruck, err := json.Marshal(newTruck)
	if err != nil {
		t.Error("Could not marshal json")
	}
	requestBody := strings.NewReader(string(jsonNewTruck))
	req, err := http.NewRequest("POST", "/trucks", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /trucks")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create truck success")
	
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	response := string(body)
	wanted := "Debe especificar la data del camión"
	if response != wanted {
		t.Errorf("response = %v, wanted %v", response, wanted)
	}
}

func TestCreateTruckRepeatedName(t *testing.T) {
	router := httprouter.New()
	router.POST("/trucks", CreateTruck)

	newTruck := types.Truck{}
	newTruck.Name = "primer camion"
	newTruck.Data = "repeated name"

	jsonNewTruck, err := json.Marshal(newTruck)
	if err != nil {
		t.Error("Could not marshal json")
	}
	requestBody := strings.NewReader(string(jsonNewTruck))
	req, err := http.NewRequest("POST", "/trucks", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /trucks")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create truck success")
	
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	response := string(body)
	wanted := "El nombre del camión ya ha sido utilizado"
	if response != wanted {
		t.Errorf("response = %v, wanted %v", response, wanted)
	}
}

func TestPatchTruck(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttruck", GetLastTruck)

	var lastTruck types.Truck

	req, err := http.NewRequest("GET", "/lasttruck", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /lasttruck")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last truck")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Truck type")
	err = json.Unmarshal(body, &lastTruck)
	if err != nil {
		t.Error("Response is not of type Truck")
	}

	router.PATCH("/trucks/:id", PatchTruck)
	requestUrl := fmt.Sprintf("/trucks/%v", lastTruck.Id)
	newTruck := types.Truck{}
	newTruck.Name = utils.RandStringBytes(10)
	newTruck.Data = "patch truck data"

	jsonNewTruck, err := json.Marshal(newTruck)
	if err != nil {
		t.Error("Could not marshal newTruck into json format")
	}
	requestBody  := strings.NewReader(string(jsonNewTruck))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /trucks/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for patching last truck")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing getting back an actor from patch request")
	responseTruck := types.Truck{}
	body2, err := ioutil.ReadAll(rr2.Body)

	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body2, &responseTruck)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Truck type")
	}
}

func TestPatchTruckEmptyName(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttruck", GetLastTruck)

	var lastTruck types.Truck

	req, err := http.NewRequest("GET", "/lasttruck", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /lasttruck")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last truck")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Truck type")
	err = json.Unmarshal(body, &lastTruck)
	if err != nil {
		t.Error("Response is not of type Truck")
	}

	router.PATCH("/trucks/:id", PatchTruck)
	requestUrl := fmt.Sprintf("/trucks/%v", lastTruck.Id)
	newTruck := types.Truck{}
	newTruck.Data = "patch truck data"

	jsonNewTruck, err := json.Marshal(newTruck)
	if err != nil {
		t.Error("Could not marshal newTruck into json format")
	}
	requestBody  := strings.NewReader(string(jsonNewTruck))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /trucks/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for patching last truck")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	
	response := string(body2)
	wanted := "Debe especificar el nombre del camión"
	if response != wanted {
		t.Errorf("response = %v, wanted %v", response, wanted)
	}
}

func TestPatchTruckEmptyDescription(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttruck", GetLastTruck)

	var lastTruck types.Truck

	req, err := http.NewRequest("GET", "/lasttruck", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /lasttruck")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last truck")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Truck type")
	err = json.Unmarshal(body, &lastTruck)
	if err != nil {
		t.Error("Response is not of type Truck")
	}

	router.PATCH("/trucks/:id", PatchTruck)
	requestUrl := fmt.Sprintf("/trucks/%v", lastTruck.Id)
	newTruck := types.Truck{}
	newTruck.Name = utils.RandStringBytes(10)

	jsonNewTruck, err := json.Marshal(newTruck)
	if err != nil {
		t.Error("Could not marshal newTruck into json format")
	}
	requestBody  := strings.NewReader(string(jsonNewTruck))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /trucks/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for patching last truck")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	
	response := string(body2)
	wanted := "Debe especificar la data del camión"
	if response != wanted {
		t.Errorf("response = %v, wanted %v", response, wanted)
	}
}

func TestPatchTruckRepeatedName(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttruck", GetLastTruck)

	var lastTruck types.Truck

	req, err := http.NewRequest("GET", "/lasttruck", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /lasttruck")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last truck")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Truck type")
	err = json.Unmarshal(body, &lastTruck)
	if err != nil {
		t.Error("Response is not of type Truck")
	}

	router.PATCH("/trucks/:id", PatchTruck)
	requestUrl := fmt.Sprintf("/trucks/%v", lastTruck.Id)
	newTruck := types.Truck{}
	newTruck.Name = "primer camion"
	newTruck.Data = "repeated name data"

	jsonNewTruck, err := json.Marshal(newTruck)
	if err != nil {
		t.Error("Could not marshal newTruck into json format")
	}
	requestBody  := strings.NewReader(string(jsonNewTruck))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /trucks/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for patching last truck")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	
	response := string(body2)
	wanted := "El nombre del camión ya ha sido utilizado"
	if response != wanted {
		t.Errorf("response = %v, wanted %v", response, wanted)
	}
}

func TestPatchTruckWrongId(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttruck", GetLastTruck)

	var lastTruck types.Truck

	req, err := http.NewRequest("GET", "/lasttruck", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /lasttruck")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last truck")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Truck type")
	err = json.Unmarshal(body, &lastTruck)
	if err != nil {
		t.Error("Response is not of type Truck")
	}

	router.PATCH("/trucks/:id", PatchTruck)
	requestUrl := fmt.Sprintf("/trucks/%v", 0)
	newTruck := types.Truck{}
	newTruck.Name = "primer camion"
	newTruck.Data = "repeated name data"

	jsonNewTruck, err := json.Marshal(newTruck)
	if err != nil {
		t.Error("Could not marshal newTruck into json format")
	}
	requestBody  := strings.NewReader(string(jsonNewTruck))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /trucks/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for patching last truck")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	
	response := string(body2)
	wanted := "Id de camión no válido"
	if response != wanted {
		t.Errorf("response = %v, wanted %v", response, wanted)
	}
}

func TestPatchTruckNonExistingId(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttruck", GetLastTruck)

	var lastTruck types.Truck

	req, err := http.NewRequest("GET", "/lasttruck", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /lasttruck")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last truck")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Truck type")
	err = json.Unmarshal(body, &lastTruck)
	if err != nil {
		t.Error("Response is not of type Truck")
	}

	router.PATCH("/trucks/:id", PatchTruck)
	requestUrl := fmt.Sprintf("/trucks/%v", 9999999)
	newTruck := types.Truck{}
	newTruck.Name = "primer camion"
	newTruck.Data = "repeated name data"

	jsonNewTruck, err := json.Marshal(newTruck)
	if err != nil {
		t.Error("Could not marshal newTruck into json format")
	}
	requestBody  := strings.NewReader(string(jsonNewTruck))
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /trucks/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for patching last truck")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	
	response := string(body2)
	wanted := "El camión especificado no existe"
	if response != wanted {
		t.Errorf("response = %v, wanted %v", response, wanted)
	}
}

func TestDeleteTruck(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttruck", GetLastTruck)

	var lastTruck types.Truck

	req, err := http.NewRequest("GET", "/lasttruck", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /lasttruck")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last truck")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Truck type")
	err = json.Unmarshal(body, &lastTruck)
	if err != nil {
		t.Error("Response is not of type Truck")
	}

	router.DELETE("/trucks/:id", DeleteTruck)
	requestUrl := fmt.Sprintf("/trucks/%v", lastTruck.Id)
	 
	req2, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /trucks/:id")
	}

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for deleting last actor")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing getting back an responseId from delete request")
	responseId := types.IdResponse{}
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body2, &responseId)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an IdResponse type")
	}
}

func TestDeleteTruckBadId(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttruck", GetLastTruck)

	var lastTruck types.Truck

	req, err := http.NewRequest("GET", "/lasttruck", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /lasttruck")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last truck")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Truck type")
	err = json.Unmarshal(body, &lastTruck)
	if err != nil {
		t.Error("Response is not of type Truck")
	}

	router.DELETE("/trucks/:id", DeleteTruck)
	requestUrl := fmt.Sprintf("/trucks/%v", 0)
	 
	req2, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /trucks/:id")
	}

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for deleting last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an responseId from delete request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	response := string(body2)
	wanted := "Id de camión no válido"
	if response != wanted {
		t.Errorf("response = %v, wanted %v", response, wanted)
	}
}

func TestDeleteTruckNonExistingId(t *testing.T) {
	router := httprouter.New()
	router.GET("/lasttruck", GetLastTruck)

	var lastTruck types.Truck

	req, err := http.NewRequest("GET", "/lasttruck", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /lasttruck")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last truck")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Truck type")
	err = json.Unmarshal(body, &lastTruck)
	if err != nil {
		t.Error("Response is not of type Truck")
	}

	router.DELETE("/trucks/:id", DeleteTruck)
	requestUrl := fmt.Sprintf("/trucks/%v", 9999)
	 
	req2, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /trucks/:id")
	}

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for deleting last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an responseId from delete request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	response := string(body2)
	wanted := "El camión con el id 9999 no existe"
	if response != wanted {
		t.Errorf("response = %v, wanted %v", response, wanted)
	}
}