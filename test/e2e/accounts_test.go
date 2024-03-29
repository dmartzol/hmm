//go:build e2e
// +build e2e

// Package httptest_test provides functional tests
package httptest_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/dghubble/sling"
)

func TestCreateAccount(t *testing.T) {
	body := strings.NewReader(`{
		"FirstName": "daniel",
		"LastName": "Martinez",
		"Email": "myemail@example.com",
		"DOB": "1990-01-01T00:00:00Z",
		"Password": "password123"
	}`)

	req, err := sling.New().
		Base("http://localhost:1100/").
		Post("accounts").
		Body(body).Request()
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error doing request: %+v", err)
	}

	if resp.StatusCode > 300 {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, resp.Body)
		t.Logf("response body: %s", buf.String())
		t.Fatalf("code: %d", resp.StatusCode)
	}

	// TODO: check if response if the expected one
}
