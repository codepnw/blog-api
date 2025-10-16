package errs

import "errors"

// Post
var (
	ErrPostNotFound = errors.New("post not found")
)

// User
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserUnauthorized = errors.New("unauthorized to update user")
	ErrUserInvalid      = errors.New("invalid email or password")
)
