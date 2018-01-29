package validation

import (
	"testing"

	"github.com/moby/tool/src/moby"
	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	minimalYaml := `
kernel:
  image: linuxkit/kernel:4.9.69
  cmdline: "console=tty0 console=ttyS0 console=ttyAMA0"
files:
  - path: /etc/issue
    contents: "wellcome to EliotOS"
trust:
  org:
    - linuxkit`
	c, parseErr := moby.NewConfig([]byte(minimalYaml))
	assert.NoError(t, parseErr)

	err := IsValid(c)
	assert.NoError(t, err)
}

func TestIsValidReturnErrorIfRelativeFiles(t *testing.T) {
	yamlWithFiles := `
kernel:
  image: linuxkit/kernel:4.9.69
  cmdline: "console=tty0 console=ttyS0 console=ttyAMA0"
files:
  - path: /etc/issue
    source: "/some/path/in/the/server"
trust:
  org:
    - linuxkit`
	c, parseErr := moby.NewConfig([]byte(yamlWithFiles))
	assert.NoError(t, parseErr)

	err := IsValid(c)
	assert.Error(t, err)
}
