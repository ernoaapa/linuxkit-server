package api

import (
	"archive/tar"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

	tw := tar.NewWriter(w)
	defer tw.Close()

	if err := dirtar.Tar(tempDir, w); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
