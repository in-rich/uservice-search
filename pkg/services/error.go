package services

import "errors"

var (
	ErrInvalidNoteUpdate     = errors.New("invalid note update")
	ErrInvalidNoteSearch     = errors.New("invalid note search")
	ErrNoteNotFound          = errors.New("note not found")
	ErrInvalidMessageUpdate  = errors.New("invalid message update")
	ErrInvalidMessageSearch  = errors.New("invalid message search")
	ErrInvalidTeamMetaCreate = errors.New("invalid team meta creation")
	ErrInvalidTeamMetaDelete = errors.New("invalid team meta deletion")
)
