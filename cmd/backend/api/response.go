package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type JSONError struct {
	Error      string
	StatusCode int
}

func (a Resources) Unmarshal(r *http.Request, iface interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll: %+v", err)
		return err
	}
	err = json.Unmarshal(body, &iface)
	if err != nil {
		log.Printf("json.Unmarshal %+v", err)
		return err
	}
	return nil
}

func (a Resources) RespondText(w http.ResponseWriter, text string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, text)
	if err != nil {
		a.Logger.Errorf("unable to write response: %v", err)
	}
}

func (a Resources) RespondJSON(w http.ResponseWriter, object interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	err := json.NewEncoder(w).Encode(object)
	if err != nil {
		a.Logger.Errorf("unable to write json response: %v", err)
	}
}

func (a Resources) RespondJSONError(w http.ResponseWriter, errorMessage string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if errorMessage == "" {
		errorMessage = http.StatusText(code)
	}
	e := JSONError{
		Error:      errorMessage,
		StatusCode: code,
	}
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		a.Logger.Errorf("unable to write json response error: %v", err)
	}
}
