package houseErrors

import "errors"

var (
	// ErrNotAuthenticated ...
	ErrNotAuthenticated = errors.New("not_authenticated")

	// ErrAlreadyExists ...
	ErrAlreadyExists = errors.New("already_exists")

	// ErrNotFound ...
	ErrNotFound = errors.New("not_found")

	// ErrWrongArgument ...
	ErrWrongArgument = errors.New("wrong_argument")

	// ErrInternalError ...
	ErrInternalError = errors.New("internal_error")

	// ErrWrongID  ...
	ErrWrongID = errors.New("wrong_id")
)
