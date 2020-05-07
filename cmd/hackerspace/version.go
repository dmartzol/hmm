package main

import (
	"fmt"
	"net/http"
)

func (api API) version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "version %s", apiVersionNumber)
}
