package actors

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

func TestGetActors(t *testing.T) {
	router := httprouter.New()
	router.GET("/actors", GetActors)

	req, err := http.NewRequest("GET", "/actors", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /actos")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing for an array of actors")
	actors := []types.Actor{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	err = json.Unmarshal(body, &actors)
	if err != nil {
		t.Error("Response body does not contain an array of type Actor")
	}
}

func TestCreateActor(t *testing.T) {
	router := httprouter.New()
	router.POST("/actors", CreateActor)

	actorName := utils.RandStringBytes(5)
	actorDescription := "Any guy"
	bodyString := fmt.Sprintf(`
		{
			"Name": "%v",
			"Description": "%v"
		}
	`, actorName, actorDescription)
	requestBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/actors", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /actors")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing Ok status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing create actor success")
	requestResponse := types.Actor{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &requestResponse)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Actor type")
	}

	if requestResponse.Name != actorName {
		t.Errorf("requestResponse.Name = %v, want %v", requestResponse.Name, actorName)
	}

	if requestResponse.Description != actorDescription {
		t.Errorf("requestResponse.Description = %v, want %v", requestResponse.Description, actorDescription)
	}
}

func TestCreateActorWithoutName(t *testing.T) {
	router := httprouter.New()
	router.POST("/actors", CreateActor)

	actorName := ""
	actorDescription := "Any guy"
	bodyString := fmt.Sprintf(`
		{
			"Name": "%v",
			"Description": "%v"
		}
	`, actorName, actorDescription)
	requestBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/actors", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /actors")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create actor success")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Debe especificar el nombre del actor"
	if string(body) != wanted {
		t.Errorf("response = %v, want %v", string(body), wanted)
	}
}

func TestCreateActorWithoutDescription(t *testing.T) {
	router := httprouter.New()
	router.POST("/actors", CreateActor)

	actorName := utils.RandStringBytes(5)
	actorDescription := ""
	bodyString := fmt.Sprintf(`
		{
			"Name": "%v",
			"Description": "%v"
		}
	`, actorName, actorDescription)
	requestBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/actors", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /actors")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create actor success")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Debe especificar la descripción del actor"
	if string(body) != wanted {
		t.Errorf("response = %v, want %v", string(body), wanted)
	}
}

func TestCreateActorWithBadJson(t *testing.T) {
	router := httprouter.New()
	router.POST("/actors", CreateActor)

	actorName := utils.RandStringBytes(5)
	actorDescription := "abc"
	bodyString := fmt.Sprintf(`
		{
			"Name": "%v",
			"Description": "%v",
		}
	`, actorName, actorDescription)
	requestBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/actors", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /actors")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create actor success")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "La data recibida no es del tipo Actor"
	if string(body) != wanted {
		t.Errorf("response = %v, want %v", string(body), wanted)
	}
}

func TestCreateActorWithRepeatedName(t *testing.T) {
	router := httprouter.New()
	router.POST("/actors", CreateActor)

	actorName := utils.RandStringBytes(15)
	actorDescription := "Any guy"
	bodyString := fmt.Sprintf(`
		{
			"Name": "%v",
			"Description": "%v"
		}
	`, actorName, actorDescription)
	requestBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/actors", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /actors")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing Ok status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	requestBody2 := strings.NewReader(bodyString)
	req2, err := http.NewRequest("POST", "/actors", requestBody2)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /actors")
	}
	rr2 := httptest.NewRecorder()

	router.ServeHTTP(rr2, req2)
	t.Log("testing bad request status code")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create actor with repeated name fail")
	body, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	wanted := "El nombre ya ha sido utilizado"
	if string(body) != wanted {
		t.Errorf("reponse = %v, wanted %v", string(body), wanted)
	}
}

// func TestPatchActor(t *testing.T) {
// 	router := httprouter.New()
// 	router.GET("/lastactorid", GetLastActorId)

// 	var lastId string

// 	req, err := http.NewRequest("GET", "/lastactorid", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Error("Could not make a get request to /lastactorid")
// 	}

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	t.Log("testing OK request status code for getting last actor id")
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("status = %v, want %v", status, http.StatusOK)
// 	}

// 	body, err := ioutil.ReadAll(rr.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 		t.Error("Could not read body of response")
// 	}

// 	lastId = string(body)

// 	router.PATCH("/actors", PatchActor)
// 	id := strconv.Atoi(lastId)
// }