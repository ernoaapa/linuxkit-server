package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigWithNormalLinuxkitConfig(t *testing.T) {
	_, err := NewConfig([]byte(`
kernel:
  image: linuxkit/kernel:4.9.69
  cmdline: "console=tty0 console=ttyS0 console=ttyAMA0"
trust:
  org:
    - linuxkit`))
	assert.NoError(t, err)
}

func TestNewConfigWithExtraFields(t *testing.T) {
	config, err := NewConfig([]byte(`
image:
  partitions:
    - { start: 0, size: 104857600, boot: true, type: "fat32" }
    - { start: 104857600, size: 209715200, type: "ext4" }
kernel:
  image: linuxkit/kernel:4.9.69
  cmdline: "console=tty0 console=ttyS0 console=ttyAMA0"
trust:
  org:
    - linuxkit`))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(config.Image.Partitions))
	assert.Equal(t, "fat32", config.Image.Partitions[0].FsType)
}
