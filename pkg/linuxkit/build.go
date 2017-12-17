package linuxkit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/moby/tool/src/moby"
	"github.com/pkg/errors"
)

func Build(name string, config []byte, formats []string, targetDir string) (err error) {
	c, err := moby.NewConfig(config)
	if err != nil {
		return NewErrInvalidConfiguration(err.Error())
	}

	var tempTar *os.File
	if tempTar, err = ioutil.TempFile("", ""); err != nil {
		return fmt.Errorf("Error creating tempfile: %v", err)
	}
	defer os.Remove(tempTar.Name())

	err = moby.Build(c, tempTar, false, "")
	if err != nil {
		return NewErrBuildFailed(err.Error())
	}

	if err := tempTar.Close(); err != nil {
		return errors.Wrap(err, "Error closing tempfile")
	}

	err = moby.Formats(filepath.Join(targetDir, name), tempTar.Name(), formats, 1024)
	if err != nil {
		return NewErrBuildFailed(err.Error())
	}
	return nil
}
