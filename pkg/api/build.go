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
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func createBuild(name, format, output string, w http.ResponseWriter, r *http.Request) {
	log.Debugf("create build, name: %s, format: %s, output: %s", name, format, output)
	body, _ := ioutil.ReadAll(r.Body)

	buildDir, err := ioutil.TempDir("", "linuxkit")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create temporary build directory: %s", err), 500)
	}
	defer os.RemoveAll(buildDir)

	if err := linuxkit.Build(name, body, []string{format}, buildDir); err != nil {
		if linuxkit.IsInvalidConfiguration(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if linuxkit.IsBuildFailed(err) {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := writeResponse(buildDir, name, format, output, w); err != nil {
		log.Debugf("Writing response caused error: %s", err)
		http.Error(w, err.Error(), 500)
	}
}

func writeResponse(buildDir, name, format, output string, w io.Writer) error {
	log.Debugf("write response buildDir: %s, name: %s, format: %s, output: %s", buildDir, name, format, output)
	switch format {
	case "rpi3":
		tar, err := os.Open(filepath.Join(buildDir, fmt.Sprintf("%s.tar", name)))
		if err != nil {
			return errors.Wrap(err, "Failed to open rpi3 tar file")
		}
		defer tar.Close()

		switch output {
		case "img":
			log.Debugf("Write img output")
			tempDir, err := ioutil.TempDir("", "img-build")
			if err != nil {
				return errors.Wrap(err, "Failed to create temporary unpacking directory")
			}
			defer os.RemoveAll(tempDir)

			if err := dirtar.Untar(tempDir, tar); err != nil {
				return errors.Wrap(err, "Error while unpacking rpi3 package")
			}

			if err := image.Build(tempDir, w); err != nil {
				return errors.Wrap(err, "Failed to build img file")
			}
		default:
			log.Debugf("Write tar output")
			_, err := io.Copy(w, tar)
			if err != nil {
				return errors.Wrap(err, "Error while copying tar file to response")
			}
		}

	default:
		if err := dirtar.Tar(buildDir, w); err != nil {
			return errors.Wrap(err, "Failed to build response tar package")
		}
	}

	return nil
}
