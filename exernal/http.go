package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func NewNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	if len(*s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func NewNullInt(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

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

func HTTPRespond(w http.ResponseWriter, text string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintf(w, text)
}
