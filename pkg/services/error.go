package services

import "errors"

var (
	ErrInvalidNoteUpdate = errors.New("invalid note update")
	ErrInvalidNoteSearch = errors.New("invalid note search")
	ErrNoteNotFound      = errors.New("not not found")
)
