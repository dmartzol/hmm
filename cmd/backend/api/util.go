package api

import (
	"fmt"
	"net/http"
	"strings"
)

func (h Resources) Version(w http.ResponseWriter, r *http.Request) {
	versionStr := fmt.Sprintf("{\"version\": \"%s\"}", apiVersionNumber)
	_, _ = w.Write([]byte(versionStr))
}

func (h Resources) NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte(`{"message":"not implemented"}`))
}

func normalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	return name
}
