package api

import (
	"fmt"
	"math/rand"
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

// RandomConfirmationCode returns a random confirmation code of length l.
func RandomConfirmationCode(l int) string {
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte(randomInt(48, 57)) // 48 to 57 is 0 to 9 in ASCII
	}
	return string(b)
}

// randomInt returns a random integer in the interval [min, max].
func randomInt(min, max int) int {
	return min + rand.Intn(max-min) + 1
}
