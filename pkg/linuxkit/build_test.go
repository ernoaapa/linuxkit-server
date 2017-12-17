package linuxkit

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/moby/tool/src/moby"
	"github.com/stretchr/testify/assert"
)

func init() {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatalf("Error while trying to create temporary directory for moby caching: %s", err)
	}

	// When we run tests in container, we cannot use default .moby -directory for caching
	moby.MobyDir = tempDir
}

const minimalYaml = `
kernel:
  image: linuxkit/kernel:4.9.69
  cmdline: "console=tty0 console=ttyS0 console=ttyAMA0"
trust:
  org:
    - linuxkit`

func TestBuild(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "linuxkit")
	defer os.RemoveAll(tempDir)

	err := Build("testing", []byte(minimalYaml), []string{"kernel+initrd"}, tempDir)
	assert.NoError(t, err)

	assert.True(t, fileExist(path.Join(tempDir, "testing-cmdline")), "Expect to create *-cmdline build output file")
	assert.True(t, fileExist(path.Join(tempDir, "testing-initrd.img")), "Expect to create *-initrd.img build output file")
	assert.True(t, fileExist(path.Join(tempDir, "testing-kernel")), "Expect to create *-kernel build output file")
}

func fileExist(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if stat.Size() == 0 {
		return false
	}
	return true
}
