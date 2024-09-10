package dao

import (
	"context"
	"errors"

	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type CreateNoteData struct {
	NoteContent string
	TargetName  string
}

type CreateNoteRepository interface {
	CreateNote(ctx context.Context, authorID string, noteID string, data *CreateNoteData) (*entities.Note, error)
}

type createNoteRepositoryImpl struct {
	db bun.IDB
}

func (r *createNoteRepositoryImpl) CreateNote(ctx context.Context, authorID string, noteID string, data *CreateNoteData) (*entities.Note, error) {
	note := &entities.Note{
		AuthorID:   authorID,
		NoteID:     noteID,
		Content:    data.NoteContent,
		TargetName: data.TargetName,
	}

	_, err := r.db.NewInsert().Model(note).Returning(searchNotesReturning).Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			return nil, ErrNoteAlreadyExists
		}

		return nil, err
	}

	return note, nil
}

func NewCreateNoteRepository(db bun.IDB) CreateNoteRepository {
	return &createNoteRepositoryImpl{
		db: db,
	}
}
