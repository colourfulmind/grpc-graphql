package ewrap

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UserIdIsRequired      = status.Error(codes.InvalidArgument, "user id is required")
	ErrInvalidEmail       = status.Error(codes.InvalidArgument, "email is invalid")
	ErrParsingRegex       = status.Error(codes.InvalidArgument, "error parsing regexp")
	ErrPasswordRequired   = status.Error(codes.InvalidArgument, "password is required")
	ErrInvalidCredentials = status.Error(codes.InvalidArgument, "incorrect email or password")
	UserAlreadyExists     = status.Error(codes.AlreadyExists, "user already exists")
	ErrConnectionTime     = status.Error(codes.DeadlineExceeded, "reached timeout while connecting to the database")
	InternalError         = status.Error(codes.Internal, "internal error")

	ErrGetComments   = status.Error(codes.Unavailable, "error getting comments")
	ErrNotAllowed    = status.Error(codes.Unavailable, "comments on post not allowed")
	ErrFoundComments = status.Error(codes.Unavailable, "comments not found")
	ErrMaxLength     = errors.New("max length exceeded")

	PostAlreadyExists = status.Error(codes.AlreadyExists, "post with the same title already exists")
	PostNotFound      = status.Error(codes.NotFound, "post not found")
	PostIDIsRequired  = status.Error(codes.InvalidArgument, "id is required")
	TitleIsRequired   = status.Error(codes.InvalidArgument, "title is required")
	TextIsRequired    = status.Error(codes.InvalidArgument, "text is required")
	ErrGetPosts       = status.Error(codes.NotFound, "error getting posts")

	EmailIsRequired    = status.Error(codes.InvalidArgument, "email is required")
	PasswordIsRequired = status.Error(codes.InvalidArgument, "password is required")

	AccessDenied = status.Error(codes.PermissionDenied, "access denied")
)
