package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dockpit/dirtar"
	"github.com/ernoaapa/linuxkit-server/pkg/linuxkit"
	log "github.com/sirupsen/logrus"
)

func createBuild(name, format string, w http.ResponseWriter, r *http.Request) {
	log.Debugf("create build, name: %s, format: %s")
	body, _ := ioutil.ReadAll(r.Body)

	tempDir, err := ioutil.TempDir("", "linuxkit")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create temporary build directory: %s", err), 500)
	}
	defer os.RemoveAll(tempDir)

	if err := linuxkit.Build(name, body, []string{format}, tempDir); err != nil {
		if linuxkit.IsInvalidConfiguration(err) {
			http.Error(w, err.Error(), 400)
		} else if linuxkit.IsBuildFailed(err) {
			http.Error(w, err.Error(), 503)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	if format == "rpi3" {
		tar, ferr := os.Open(filepath.Join(tempDir, fmt.Sprintf("%s.tar", name)))
		if ferr != nil {
			http.Error(w, fmt.Sprintf("Failed to open rpi3 tar file: %s", ferr), 500)
			return
		}
		_, err := io.Copy(w, tar)
		if err != nil {
			log.Errorf("Error while copying tar file to response: %s", err)
		}
		return
	}

	if err := dirtar.Tar(tempDir, w); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
