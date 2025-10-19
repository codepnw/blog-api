package errs

import "errors"

// Post
var (
	ErrPostNotFound = errors.New("post not found")
)

// User
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserUnauthorized = errors.New("unauthorized")
	ErrUserInvalid      = errors.New("invalid email or password")
)

// Comment
var (
	ErrCommentNotFound   = errors.New("comment not found")
	ErrCommentIsRequired = errors.New("comment is required")
	ErrCommentNotOwner   = errors.New("user not owner of comment")
)
