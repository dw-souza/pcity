package apperrors

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrIndexingInProgress    = errors.New("indexing in progress")
	ErrAmenityAlreadyExists  = errors.New("amenity already exists")
	ErrAlreadyConfirmed      = errors.New("already confirmed")
	ErrCannotConfirmOwnReport = errors.New("cannot confirm own report")
	ErrInvalidInput          = errors.New("invalid input")
)
