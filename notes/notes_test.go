package notes

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

func TestGetNotes(t *testing.T) {
	router := httprouter.New()
	router.GET("/notes", GetNotes)

	req, err := http.NewRequest("GET", "/notes", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could no make a get request to /notes")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing for an array of notes")
	notes := []types.Note{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &notes)
	if err != nil {
		t.Error("Response body does not contain an array of type Note")
	}
}

func TestCreateNote(t *testing.T) {
	router := httprouter.New()
	router.POST("/notes", CreateNote)

	description := "description create"
	urgency := "low"
	bodyString := fmt.Sprintf(`
		{
			"Description": "%v",
			"Urgency": "%v"
		}
	`, description, urgency)
	requestBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/notes", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /notes")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	t.Log("testing Ok status code")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing create note success")
	requestResponse := types.Note{}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body, &requestResponse)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Note type")
	}
	if requestResponse.Description != description {
		t.Errorf("requestResponse.Description = %v, want %v", requestResponse.Description, description)
	}
	if requestResponse.Urgency != urgency {
		t.Errorf("requestResponse.Urgency = %v, want %v", requestResponse.Urgency, urgency)
	}
}

func TestCreateNoteWithoutDescription(t *testing.T) {
	router := httprouter.New()
	router.POST("/notes", CreateNote)

	description := ""
	urgency := "low"
	bodyString := fmt.Sprintf(`
		{
			"Description": "%v",
			"Urgency": "%v"
		}
	`, description, urgency)
	requestBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/notes", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /notes")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing bad request status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create note success")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Debe especificar la descripción de la nota"
	if string(body) != wanted {
		t.Errorf("response = %v, want %v", string(body), wanted)
	}
}

func TestCreateNoteWithoutUrgency(t *testing.T) {
	router := httprouter.New()
	router.POST("/notes", CreateNote)

	description := "a"
	urgency := ""
	bodyString := fmt.Sprintf(`
		{
			"Description": "%v",
			"Urgency": "%v"
		}
	`, description, urgency)
	requestBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/notes", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /notes")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing Ok status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create note success")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Debe especificar la urgencia de la nota, la cual puede ser baja, media, alta o crítica"
	if string(body) != wanted {
		t.Errorf("response = %v, want %v", string(body), wanted)
	}
}

func TestCreateNoteBadJson(t *testing.T) {
	router := httprouter.New()
	router.POST("/notes", CreateNote)

	description := "a"
	urgency := "low"
	bodyString := fmt.Sprintf(`
		{
			"Description": "%v",
			"Urgency": "%v",
		}
	`, description, urgency)
	requestBody := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", "/notes", requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a post request to /notes")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing Ok status code")
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing create note success")
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "La data recibida no es del tipo Nota"
	if string(body) != wanted {
		t.Errorf("response = %v, want %v", string(body), wanted)
	}
}

func TestPatchNote(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastnoteid", GetLastNoteId)

	var lastNoteId types.IdResponse
	req, err := http.NewRequest("GET", "/lastnoteid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastnoteid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las note id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastNoteId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PATCH("/notes/:id", PatchNote)
	requestUrl := fmt.Sprintf("/notes/%v", lastNoteId.Id)

	description := "patch description"
	urgency := "medium"
	bodyString := fmt.Sprintf(`
	{
		"Description": "%v",
		"Urgency": "%v"
	}
	`, description, urgency)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /notes/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for patching last note")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing getting back a note from patch request")
	responnseNote := types.Note{}
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body2, &responnseNote)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Note type")
	}
}

func TestPatchNoteBadJson(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastnoteid", GetLastNoteId)

	var lastNoteId types.IdResponse
	req, err := http.NewRequest("GET", "/lastnoteid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastnoteid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las note id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastNoteId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PATCH("/notes/:id", PatchNote)
	requestUrl := fmt.Sprintf("/notes/%v", lastNoteId.Id)

	description := "patch description"
	urgency := "medium"
	bodyString := fmt.Sprintf(`
	{
		"Description": "%v",
		"Urgency": "%v",
	}
	`, description, urgency)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /notes/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a note from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "La data enviada con corresponde a una nota"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestPatchNoteEmptyDescription(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastnoteid", GetLastNoteId)

	var lastNoteId types.IdResponse
	req, err := http.NewRequest("GET", "/lastnoteid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastnoteid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las note id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastNoteId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PATCH("/notes/:id", PatchNote)
	requestUrl := fmt.Sprintf("/notes/%v", lastNoteId.Id)

	description := ""
	urgency := "medium"
	bodyString := fmt.Sprintf(`
	{
		"Description": "%v",
		"Urgency": "%v"
	}
	`, description, urgency)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /notes/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a note from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Debe especificar la descripción de la nota que desea modificar"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestPatchNoteBadUrgency(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastnoteid", GetLastNoteId)

	var lastNoteId types.IdResponse
	req, err := http.NewRequest("GET", "/lastnoteid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastnoteid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las note id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastNoteId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PATCH("/notes/:id", PatchNote)
	requestUrl := fmt.Sprintf("/notes/%v", lastNoteId.Id)

	description := "patch description"
	urgency := "bad"
	bodyString := fmt.Sprintf(`
	{
		"Description": "%v",
		"Urgency": "%v"
	}
	`, description, urgency)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /notes/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a note from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Debe especificar la urgencia de la nota, la cual puede ser baja, media, alta o crítica"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestPatchNoteNonExistingNote(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastnoteid", GetLastNoteId)

	var lastNoteId types.IdResponse
	req, err := http.NewRequest("GET", "/lastnoteid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastnoteid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las note id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastNoteId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PATCH("/notes/:id", PatchNote)
	requestUrl := fmt.Sprintf("/notes/%v", lastNoteId.Id)

	description := "patch description"
	urgency := "bad"
	bodyString := fmt.Sprintf(`
	{
		"Description": "%v",
		"Urgency": "%v"
	}
	`, description, urgency)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /notes/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a note from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Debe especificar la urgencia de la nota, la cual puede ser baja, media, alta o crítica"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestPatchNoteWrongId(t *testing.T) {
	router := httprouter.New()

	router.PATCH("/notes/:id", PatchNote)
	requestUrl := fmt.Sprintf("/notes/%v", 0)

	description := "patch description"
	urgency := "bad"
	bodyString := fmt.Sprintf(`
	{
		"Description": "%v",
		"Urgency": "%v"
	}
	`, description, urgency)
	requestBody := strings.NewReader(bodyString)
	req2, err := http.NewRequest("PATCH", requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a patch request to /notes/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for patching last note")
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", status, http.StatusBadRequest)
	}

	t.Log("testing getting back a note from patch request")
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	wanted := "Id de nota no válido"
	if string(body2) != wanted {
		t.Errorf("response = %v, wanted %v", string(body2), wanted)
	}
}

func TestDeleteActor(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastnoteid", GetLastNoteId)

	var lastNoteId types.IdResponse
	req, err := http.NewRequest("GET", "/lastnoteid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastnoteid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las note id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastNoteId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.DELETE("/notes/:id", DeleteNote)
	requestUrl := fmt.Sprintf("/notes/%v", lastNoteId.Id)

	req2, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /notes/:id")
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
	if response.Id != lastNoteId.Id {
		t.Errorf("response.Id = '%v', lastNoteId.Id = '%v', they are different", response.Id, lastNoteId.Id)
	}
}

func TestDeleteActorZeroId(t *testing.T) {
	router := httprouter.New()

	router.DELETE("/notes/:id", DeleteNote)
	requestUrl := fmt.Sprintf("/notes/%v", 0)

	req2, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /notes/:id")
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
	wanted := "Id de nota no válido"
	if wanted != string(body2) {
		t.Errorf("Response body = '%v', wanted = '%v'", string(body2), wanted)
	}
}

func TestDeleteActorbadId(t *testing.T) {
	router := httprouter.New()

	router.DELETE("/notes/:id", DeleteNote)
	requestUrl := fmt.Sprintf("/notes/%v", 9999)

	req2, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /notes/:id")
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
	wanted := "La nota con el id 9999 no existe"
	if wanted != string(body2) {
		t.Errorf("Response body = '%v', wanted = '%v'", string(body2), wanted)
	}
}

func TestAttendNote(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastnoteid", GetLastNoteId)

	var lastNoteId types.IdResponse
	req, err := http.NewRequest("GET", "/lastnoteid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastnoteid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las note id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastNoteId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PUT("/attend_note/:id", AttendNote)
	requestUrl := fmt.Sprintf("/attend_note/%v", lastNoteId.Id)

	req2, err := http.NewRequest("PUT", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /notes/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for deleting last note")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing getting back an IdResponse from delete request")
	response := types.Note{}
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body2, &response)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Note type")
	}

	t.Log("testing response.Attended is true")
	if !response.Attended  {
		t.Errorf("response.Attended = '%v', want true", response.Attended)
	}
}

func TestUnattendNote(t *testing.T) {
	router := httprouter.New()
	router.GET("/lastnoteid", GetLastNoteId)

	var lastNoteId types.IdResponse
	req, err := http.NewRequest("GET", "/lastnoteid", nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a get request to /lastnoteid")
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Log("testing OK request status code for getting las note id")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}

	t.Log("testing body response to be of IdResponse type")
	err = json.Unmarshal(body, &lastNoteId)
	if err != nil {
		t.Error("Response is not of type IdResponse")
	}

	router.PUT("/unattend_note/:id", UnattendNote)
	requestUrl := fmt.Sprintf("/unattend_note/%v", lastNoteId.Id)

	req2, err := http.NewRequest("PUT", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not make a delete request to /notes/:id")
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	t.Log("testing OK request status code for deleting last note")
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("status = %v, want %v", status, http.StatusOK)
	}

	t.Log("testing getting back an IdResponse from delete request")
	response := types.Note{}
	body2, err := ioutil.ReadAll(rr2.Body)
	if err != nil {
		log.Fatal(err)
		t.Error("Could not read body of response")
	}
	err = json.Unmarshal(body2, &response)
	if err != nil {
		log.Fatal(err)
		t.Error("Response body does not contain an Note type")
	}

	t.Log("testing response.Attended is true")
	if response.Attended  {
		t.Errorf("response.Attended = '%v', want false", response.Attended)
	}
}