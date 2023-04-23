package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

// Version returns the API version.
func (h Resources) Version(w http.ResponseWriter, r *http.Request) {
	versionStr := fmt.Sprintf("{\"version\": \"%s\"}", apiVersionNumber)
	_, _ = w.Write([]byte(versionStr))
}

// NotImplementedHandler responds with a 501 status code and a message indicating that the requested functionality has not been implemented.
func (h Resources) NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte(`{"message":"not implemented"}`))
}

// normalizeName normalizes a string by removing leading and trailing white spaces and
// converting all characters to lowercase.
func normalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	return name
}

// RandomConfirmationCode generates a random confirmation code of a specified length
// using digits 0-9.
func RandomConfirmationCode(length int) string {
	chars := []rune("0123456789")
	return randomSample(chars, length)
}

// RandomPassword generates a random password of a specified length using
// uppercase and lowercase letters, numbers, and special characters.
func RandomPassword(length int) string {
	letters := []rune(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)
	numbers := []rune(`0123456789`)
	specials := []rune(`!@#$%&*-+=.?`)
	chars := append(letters, numbers...)
	chars = append(chars, specials...)
	return randomSample(chars, length)
}

// randomSample generates a random string of a specified length using characters
// from a provided character set.
func randomSample(charSet []rune, length int) string {
	var builder strings.Builder
	for i := 0; i < length; i++ {
		randomChar := charSet[rand.Intn(len(charSet))]
		builder.WriteRune(randomChar)
	}
	return builder.String()
}
