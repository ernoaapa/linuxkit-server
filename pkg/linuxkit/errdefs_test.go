package linuxkit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInvalidConfiguration(t *testing.T) {
	assert.True(t, IsInvalidConfiguration(NewErrInvalidConfiguration("error message")))
	assert.False(t, IsInvalidConfiguration(NewErrBuildFailed("error message")))
}

func TestBuildFailed(t *testing.T) {
	assert.True(t, IsBuildFailed(NewErrBuildFailed("error message")))
	assert.False(t, IsBuildFailed(NewErrInvalidConfiguration("error message")))
}
