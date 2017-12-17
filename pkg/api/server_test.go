package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

const minimalYaml = `
kernel:
  image: linuxkit/kernel:4.9.69
  cmdline: "console=tty0 console=ttyS0 console=ttyAMA0"
trust:
  org:
    - linuxkit`

func TestCreateBuild(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		createBuild("testing", "kernel+initrd", w, r)
	}))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(minimalYaml)))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}
