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

func TestGetCompanies(t *testing.T) {
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

	actorType := "third"
	actorName := utils.RandStringBytes(5)
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
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
}

func TestCreateActorWithoutName(t *testing.T) {
	router := httprouter.New()
	router.POST("/actors", CreateActor)

	actorType := "mine"
	actorName := ""
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
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

func TestCreateActorWithBadJson(t *testing.T) {
	router := httprouter.New()
	router.POST("/actors", CreateActor)

	actorType := "personnel"
	actorName := utils.RandStringBytes(20)
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v",
		}
	`, actorType, actorName)
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

	actorType := "third"
	actorName := "Externo"
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
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

func TestCreateActorWithBadType(t *testing.T) {
	router := httprouter.New()
	router.POST("/actors", CreateActor)

	actorType := "notype"
	actorName := "Externo"
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
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
	wanted := "Debe especificar el tipo de actor, el cual puede ser 'personal', 'tercero', 'mina' o 'contratante'"
	if string(body) != wanted {
		t.Errorf("reponse = %v, wanted %v", string(body), wanted)
	}
}

func TestPatchActor(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastactor", GetLastActor)

	var lastActor types.Actor

	req, err := http.NewRequest("GET", "/lastactor", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastactor")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last actor")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Actor type")
	err = json.Unmarshal(body, &lastActor)
	if err != nil {
		t.Error("Response is not of type Actor")
	}

	router.PATCH("/actors/:id", PatchActor)
	requestUrl := fmt.Sprintf("/actors/%v", lastActor.Id)

	actorType := "contractee"
	actorName := utils.RandStringBytes(20)
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for patching last actor")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing getting back an actor from patch request")
	responnseActor := types.Actor{}
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body2, &responnseActor)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Actor type")
	}
}

func TestPatchActorBadJson(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastactor", GetLastActor)

	var lastActor types.Actor

	req, err := http.NewRequest("GET", "/lastactor", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastactor")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last actor")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Actor type")
	err = json.Unmarshal(body, &lastActor)
	if err != nil {
		t.Error("Response is not of type Actor")
	}

	router.PATCH("/actors/:id", PatchActor)
	requesUrl := fmt.Sprintf("/actors/%v", lastActor.Id)

	actorType := "mine"
	actorName := utils.RandStringBytes(20)
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v",
		}
	`, actorType, actorName)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requesUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing bad request status code for patching last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "La data enviada con corresponde con un actor parcial"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestPatchActorEmptyName(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastactor", GetLastActor)

	var lastActor types.Actor

	req, err := http.NewRequest("GET", "/lastactor", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastactor")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last actor")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Actor type")
	err = json.Unmarshal(body, &lastActor)
	if err != nil {
		t.Error("Response is not of type Actor")
	}

	router.PATCH("/actors/:id", PatchActor)
	requesUrl := fmt.Sprintf("/actors/%v", lastActor.Id)

	actorType := "personnel"
	actorName := ""
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requesUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing bad request status code for patching last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Debe especificar el nombre del actor que desea modificar"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestPatchActorDuplicatedName(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastactor", GetLastActor)

	var lastActor types.Actor

	req, err := http.NewRequest("GET", "/lastactor", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastactor")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last actor")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Actor type")
	err = json.Unmarshal(body, &lastActor)
	if err != nil {
		t.Error("Response is not of type Actor")
	}

	router.PATCH("/actors/:id", PatchActor)
	requesUrl := fmt.Sprintf("/actors/%v", lastActor.Id)

	actorType := "contractee"
	actorName := "Externo"
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requesUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing bad request status code for patching last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "El nombre ya ha sido utilizado"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestPatchActorNonExistingActor(t *testing.T) {
	router := httprouter.New()

	router.PATCH("/actors/:id", PatchActor)
	requesUrl := fmt.Sprintf("/actors/%v", 9999)

	actorType := "third"
	actorName := utils.RandStringBytes(10)
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requesUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing bad request status code for patching last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "El actor especificado no existe"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestPatchActorExterno(t *testing.T) {
	router := httprouter.New()

	router.PATCH("/actors/:id", PatchActor)
	requesUrl := fmt.Sprintf("/actors/%v", 1)

	actorType := "mine"
	actorName := utils.RandStringBytes(10)
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requesUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing bad request status code for patching last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "No puede modificar el actor externo"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestPatchActorWrongId(t *testing.T) {
	router := httprouter.New()

	router.PATCH("/actors/:id", PatchActor)
	requesUrl := fmt.Sprintf("/actors/%v", 0)

	actorType := "personnel"
	actorName := utils.RandStringBytes(10)
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requesUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing bad request status code for patching last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Id de actor no válido"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestPatchActorWrongType(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastactor", GetLastActor)

	var lastActor types.Actor

	req, err := http.NewRequest("GET", "/lastactor", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastactor")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last actor")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Actor type")
	err = json.Unmarshal(body, &lastActor)
	if err != nil {
		t.Error("Response is not of type Actor")
	}

	router.PATCH("/actors/:id", PatchActor)
	requesUrl := fmt.Sprintf("/actors/%v", lastActor.Id)

	actorType := "notype"
	actorName := ""
	bodyString := fmt.Sprintf(`
		{
			"Type": "%v",
			"Name": "%v"
		}
	`, actorType, actorName)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requesUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing bad request status code for patching last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Debe especificar el tipo de actor, el cual puede ser 'personal', 'tercero', 'mina' o 'contratante'"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestDeleteActor(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastactor", GetLastActor)

	var lastActor types.Actor

	req, err := http.NewRequest("GET", "/lastactor", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastactor")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting last actor")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of Actor type")
	err = json.Unmarshal(body, &lastActor)
	if err != nil {
		t.Error("Response is not of type Actor")
	}

	router.DELETE("/actors/:id", DeleteActor)
	requesUrl := fmt.Sprintf("/actors/%v", lastActor.Id)

	req2, err := http.NewRequest("DELETE", requesUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing OK request status code for patching last actor")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing getting back an actor from patch request")
	responnseActor := types.Actor{}
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body2, &responnseActor)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Actor type")
	}
}

func TestDeleteActorBadId(t *testing.T) {
	router := httprouter.New()

	router.DELETE("/actors/:id", DeleteActor)
	requesUrl := fmt.Sprintf("/actors/%v", 0)

	req2, err := http.NewRequest("DELETE", requesUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing bad request status code for patching last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Id de actor no válido"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestDeleteActorTakenActor(t *testing.T) {
	router := httprouter.New()

	router.DELETE("/actors/:id", DeleteActor)
	requesUrl := fmt.Sprintf("/actors/%v", 1)

	req2, err := http.NewRequest("DELETE", requesUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /actors/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	t.Log("testing bad request status code for patching last actor")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back an actor from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "El actor que intenta borrar tiene una o mas facturas asociadas por lo que no puede ser eliminado"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}
