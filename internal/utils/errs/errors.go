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
)
