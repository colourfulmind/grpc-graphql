package storage

import "errors"

var (
	ErrUserExists     = errors.New("user already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrConnectionTime = errors.New("connection time is expired")

	ErrGetComments   = errors.New("failed get comments")
	ErrFoundComments = errors.New("comments not found")

	ErrPostExists   = errors.New("post with the same title already exists")
	ErrPostNotFound = errors.New("post does not exist")
	ErrGetPosts     = errors.New("failed to get posts")

	ErrQuery      = errors.New("failed query")
	ErrNotAllowed = errors.New("comments on post not allowed")

	// ErrAppNotFound  = errors.New("app not found")
)
