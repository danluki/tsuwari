package seventveventapi

import "errors"

var (
	ErrUnsupportedOperation     = errors.New("Unsupported operation")
	ErrIncorrectPayload         = errors.New("Incorrect payload")
	ErrWrongSessionID           = errors.New("Wrong session id from 7TV")
	ErrNotConnectedOnRetry      = errors.New("Not connected before, while trying to retry")
	ErrLibraryBadImplementation = errors.New("7TV reported, that client send something bad")
	ErrRateLimited              = errors.New(
		"Got rate limited from 7TV, this is probably because wrong implementation",
	)
	ErrMaintenance           = errors.New("7TV maintenance")
	ErrAlreadySubscribed     = errors.New("Already subscribed")
	ErrNotSubscribed         = errors.New("Not subscribed")
	ErrInsufficientPrivilege = errors.New("Insufficient privelege")
)
