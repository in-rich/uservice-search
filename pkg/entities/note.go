package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Note struct {
	bun.BaseModel `bun:"table:notes"`

	ID *uuid.UUID `bun:"id,pk,type:uuid"`

	AuthorID string `bun:"author_id"`
	NoteID   string `bun:"note_id"`

	Content    string `bun:"content"`
	TargetName string `bun:"target_name"`
}
