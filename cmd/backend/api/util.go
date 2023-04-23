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
	return randomSample(chars, l)
}

func RandomPassword(l int) string {
	letters := []rune(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)
	numbers := []rune(`0123456789`)
	specials := []rune(`!@#$%&*-+=.?`)
	chars := append(letters, numbers...)
	chars = append(chars, specials...)
	return randomSample(chars, l)
}

// randomSample returns a random string of the given length from the given character set.
func randomSample(charSet []rune, length int) string {
	var builder strings.Builder
	for i := 0; i < length; i++ {
		randomChar := charSet[rand.Intn(len(charSet))]
		builder.WriteRune(randomChar)
	}
	return builder.String()
}
