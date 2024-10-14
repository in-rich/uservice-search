package dao

import "errors"

var (
	ErrNoteAlreadyExists     = errors.New("note already exists")
	ErrNoteNotFound          = errors.New("note not found")
	ErrMessageAlreadyExists  = errors.New("message already exists")
	ErrMessageNotFound       = errors.New("message not found")
	ErrTeamMetaAlreadyExists = errors.New("team meta already exists")
	ErrTeamMetaNotFound      = errors.New("team meta not found")
)
