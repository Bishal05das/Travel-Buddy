package domain

import "errors"

var (
	ErrImageTargetNotFound = errors.New("agency or tour not found")
	ErrImageAccessDenied   = errors.New("not allowed to change this image")
	// ErrInvalidCredentials is returned for any failed login, whether the
	// email is unknown or the password is wrong, so callers cannot probe
	// which accounts exist.
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrEmailTaken             = errors.New("email already exists")
	ErrMemberNotFound         = errors.New("member not found")
	ErrOwnerCreationForbidden = errors.New("only a platform admin can create an agency owner")
	ErrOwnerProtected         = errors.New("agency owners cannot be removed or have their access changed")
)
