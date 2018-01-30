package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dockpit/dirtar"
	"github.com/ernoaapa/linuxkit-server/pkg/image"
	"github.com/ernoaapa/linuxkit-server/pkg/linuxkit"
	log "github.com/sirupsen/logrus"
)

func createBuild(name, format, output string, w http.ResponseWriter, r *http.Request) {
	log.Debugf("create build, name: %s, format: %s")
	body, _ := ioutil.ReadAll(r.Body)

	buildDir, err := ioutil.TempDir("", "linuxkit")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create temporary build directory: %s", err), 500)
	}
	defer os.RemoveAll(buildDir)

	if err := linuxkit.Build(name, body, []string{format}, buildDir); err != nil {
		if linuxkit.IsInvalidConfiguration(err) {
			http.Error(w, err.Error(), 400)
		} else if linuxkit.IsBuildFailed(err) {
			http.Error(w, err.Error(), 503)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	switch format {
	case "rpi3":
		tar, ferr := os.Open(filepath.Join(buildDir, fmt.Sprintf("%s.tar", name)))
		if ferr != nil {
			http.Error(w, fmt.Sprintf("Failed to open rpi3 tar file: %s", ferr), 500)
			return
		}
		defer tar.Close()

		switch output {
		case "img":

			tempDir, err := ioutil.TempDir("", "img-build")
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to create temporary unpacking directory: %s", err), 500)
			}
			defer os.RemoveAll(tempDir)

			if err := dirtar.Untar(tempDir, tar); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if err := image.Build(tempDir, w); err != nil {
				http.Error(w, fmt.Sprintf("Failed to build img file: %s", err), 500)
			}

		default:
			_, err := io.Copy(w, tar)
			if err != nil {
				log.Errorf("Error while copying tar file to response: %s", err)
			}
			return
		}

	default:
		if err := dirtar.Tar(buildDir, w); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
}
