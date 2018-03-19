package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKpartxOutput(t *testing.T) {
	assert.Equal(t, []string{"/dev/mapper/loop0p1"}, mustParseDevMappings(`
loop0p1 : 0 202752 /dev/loop0 2048
loop deleted : /dev/loop0
`))

	assert.Equal(t, []string{"/dev/mapper/loop1p1", "/dev/mapper/loop1p2"}, mustParseDevMappings(`
loop1p1 : 0 202752 /dev/loop1 2048
loop1p2 : 0 202752 /dev/loop1 102400
loop deleted : /dev/loop1
`))

	assert.Equal(t, []string{"/dev/mapper/loop1p1", "/dev/mapper/loop1p2"}, mustParseDevMappings(`
loop1p2 : 0 202752 /dev/loop1 102400
loop1p1 : 0 202752 /dev/loop1 2048
loop deleted : /dev/loop1
`))
}
