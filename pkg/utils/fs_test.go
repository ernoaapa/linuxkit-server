package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDirSize(t *testing.T) {
	dir, err := os.Getwd()
	assert.NoError(t, err, "Failed to resolve current working directory for test")

	size, err := GetDirSize(dir)
	assert.NoError(t, err)
	assert.True(t, size > 0)
}
