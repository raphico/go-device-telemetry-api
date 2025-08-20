package device

import "errors"

var (
	ErrDeviceNotFound    = errors.New("device not found")
	ErrNameRequired      = errors.New("device name is required")
	ErrNameTooShort      = errors.New("device name must be at least 3 characters")
	ErrNameTooLong       = errors.New("device name must be at most 50 characters")
	ErrNameInvalidChars  = errors.New("device name may only contain letters, numbers, underscores, periods, or hyphens")
	ErrInvalidStatus     = errors.New("device status must be either 'offline' or 'online'")
	ErrInvalidDeviceType = errors.New("device type is invalid")
	ErrInvalidMetadata   = errors.New("device metadata is invalid")
)
