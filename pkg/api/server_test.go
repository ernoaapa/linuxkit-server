package api

import (
	"bytes"
	"fmt"
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
	t.Skip("httptest.Server closes connection after 5s and hangs the test for unknown reason")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		output, err := getOutputFormat(r)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		createBuild("testing", "kernel+initrd", output, w, r)
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

func TestCreateImgBuild(t *testing.T) {
	t.Skip("httptest.Server closes connection after 5s and hangs the test for unknown reason")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		output, err := getOutputFormat(r)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		if output != "img" {
			http.Error(w, "Should give 'img' output", 400)
		}
		createBuild("testing", "kernel+initrd", output, w, r)
	}))
	defer server.Close()

	resp, err := http.Post(fmt.Sprintf("%s?output=%s", server.URL, "img"), "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(minimalYaml)))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}
