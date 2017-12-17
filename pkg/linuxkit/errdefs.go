package linuxkit

// ErrInvalidConfiguration is some error in linuxkit configuration
type ErrInvalidConfiguration struct {
	message string
}

// NewErrInvalidConfiguration create new ErrInvalidConfiguration with message
func NewErrInvalidConfiguration(message string) *ErrInvalidConfiguration {
	return &ErrInvalidConfiguration{
		message: message,
	}
}

func (e *ErrInvalidConfiguration) Error() string {
	return e.message
}

// IsInvalidConfiguration return true if the error is ErrInvalidConfiguration
func IsInvalidConfiguration(err error) bool {
	switch err.(type) {
	case *ErrInvalidConfiguration:
		return true
	default:
		return false
	}
}

// ErrBuildFailed is some error in linuxkit configuration
type ErrBuildFailed struct {
	message string
}

// NewErrBuildFailed create new ErrBuildFailed with message
func NewErrBuildFailed(message string) *ErrBuildFailed {
	return &ErrBuildFailed{
		message: message,
	}
}

func (e *ErrBuildFailed) Error() string {
	return e.message
}

// IsBuildFailed return true if the error is ErrBuildFailed
func IsBuildFailed(err error) bool {
	switch err.(type) {
	case *ErrBuildFailed:
		return true
	default:
		return false
	}
}
