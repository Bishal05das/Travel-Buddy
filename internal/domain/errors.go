package domain

import "errors"

var (
	// ErrInvalidCredentials is returned for any failed login, whether the
	// email is unknown or the password is wrong, so callers cannot probe
	// which accounts exist.
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already exists")
)
