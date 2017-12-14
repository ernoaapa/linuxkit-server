package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateBuild(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		createBuild(w, r)
	}))
	defer server.Close()

	json := []byte("{}")
	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}
