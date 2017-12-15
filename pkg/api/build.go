package api

import (
	"encoding/json"
	"net/http"

	"github.com/moby/tool/src/moby"
)

func createBuild(w http.ResponseWriter, r *http.Request) {
	var build moby.Moby
	json.NewEncoder(w).Encode(build)
}
