//go:build e2e
// +build e2e

// Package httptest_test provides functional tests
package httptest_test

import (
	"net/http"
	"testing"

	"github.com/dghubble/sling"
	"github.com/dmartzol/goapi/internal/api"
)

func TestCreateAccount(t *testing.T) {
	body := &api.CreateAccountRequest{
		FirstName: "daniel",
		LastName:  "Martinez",
		Email:     "dani@example.com",
	}
	req, err := sling.New().Base("http://localhost:1100/v1/").Post("accounts").BodyJSON(body).Request()
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error doing request: %+v", err)
	}

	if resp.StatusCode > 300 {
		t.Fatalf("code: %d", resp.StatusCode)
	}
}
