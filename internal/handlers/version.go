package handlers

import (
	"fmt"
	"net/http"
)

func (api API) Version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "version %s", apiVersionNumber)
}
