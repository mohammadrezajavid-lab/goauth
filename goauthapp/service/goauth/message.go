package goauth

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrOTPNotFound       = errors.New("otp not found or expired")
	ErrInvalidOTPCode    = errors.New("invalid OTP code")
	ErrInternalService   = errors.New("internal auth service error")
	ErrUserAlreadyExists = errors.New("user already exists")
)
