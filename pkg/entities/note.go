package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const searchNotesReturning = "author_id, note_id, id"

type Note struct {
	bun.BaseModel `bun:"table:notes"`

	ID *uuid.UUID `bun:"id,pk,type:uuid"`

	AuthorID string `bun:"author_id"`
	NoteID   string `bun:"note_id"`

	Content    string `bun:"content"`
	TargetName string `bun:"target_name"`
}

func (n *Note) BeforeCreate(query *bun.InsertQuery) *bun.InsertQuery {
	return query.
		Value("content", Vectorize("A"), n.Content).
		Value("target_name", Vectorize("A"), n.TargetName).
		Returning(searchNotesReturning)
}

func (n *Note) BeforeUpdate(query *bun.UpdateQuery) *bun.UpdateQuery {
	return query.
		Column("content", "target_name").
		Value("content", Vectorize("A"), n.Content).
		Value("target_name", Vectorize("A"), n.TargetName).
		Returning(searchNotesReturning)
}
