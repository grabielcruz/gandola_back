package transactions

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
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
		t.Errorf("Wrong status")
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
	}
	want := "Server working"
	if string(body) != want {
		t.Errorf("body = %v; want %v", string(body), want)
	}
}