package api

import (
	"encoding/json"
	"net/http"
)

func createBuild(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		ok string
	}{
		ok: "true",
	})
}
