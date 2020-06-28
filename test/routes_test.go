package test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmartzol/hmmm/internal/controllers"
	"github.com/dmartzol/hmmm/internal/storage/mockdb"
)

var api *controllers.API

func init() {
	db, err := mockdb.NewMockDB()
	if err != nil {
		log.Fatal(err)
	}
	api, err = controllers.NewAPI(db)
	if err != nil {
		log.Fatalf("error starting api: %+v", err)
	}
}

func TestCreateAccount(t *testing.T) {
	var jsonStr = []byte(`{"FirstName":"Daniel","LastName":"Martinez","DOB":"2020-01-01","Gender":"M","PhoneNumber":"+14443332222","Email":"random@email.com","Password":"randompass"}`)

	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.CreateAccount)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	log.Print(rr.Body.String())
}
