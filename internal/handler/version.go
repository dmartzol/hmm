package handler

import (
	"fmt"
	"net/http"
)

func (h Handler) Version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "version %s", apiVersionNumber)
}
