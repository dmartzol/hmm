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
	chars := []rune("0123456789")
	return randomStringFromRunes(chars, l)
}

func RandomPassword(l int) string {
	letters := []rune(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)
	numbers := []rune(`0123456789`)
	specials := []rune(`!@#$%&*-+=.?`)
	chars := append(letters, numbers...)
	chars = append(chars, specials...)
	return randomStringFromRunes(chars, l)
}

func randomStringFromRunes(popuation []rune, l int) string {
	var b strings.Builder
	for i := 0; i < l; i++ {
		ch := popuation[rand.Intn(len(popuation))]
		b.WriteRune(ch)
	}
	return b.String()
}
