package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {

	app := New()
	srvr := httptest.NewTLSServer(app.Handlers())

	resp, err := srvr.Client().Get(srvr.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		t.Fatal("/health is not defined")
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("wrong status code. expected: %d, got: %d", http.StatusOK, resp.StatusCode)
	}
}
