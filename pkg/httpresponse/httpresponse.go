package httpresponse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Unmarshal(r *http.Request, iface interface{}) error {
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

func RespondText(w http.ResponseWriter, text string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintf(w, text)
}

func RespondJSON(w http.ResponseWriter, object interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	json.NewEncoder(w).Encode(object)
}
