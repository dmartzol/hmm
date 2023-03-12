package api

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
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

func randomCode(max int) (string, error) {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		log.Printf("ReadAtLeast: %+v", err)
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}
