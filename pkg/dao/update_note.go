package dao

import (
	"context"

	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type UpdateNoteData struct {
	NoteContent string
	TargetName  string
}

type UpdateNoteRepository interface {
	UpdateNote(ctx context.Context, authorID string, noteID string, data *UpdateNoteData) (*entities.Note, error)
}

type updateNoteRepositoryImpl struct {
	db bun.IDB
}

func (r *updateNoteRepositoryImpl) UpdateNote(ctx context.Context, authorID string, noteID string, data *UpdateNoteData) (*entities.Note, error) {
	note := &entities.Note{
		Content:    data.NoteContent,
		TargetName: data.TargetName,
	}

	res, err := note.BeforeUpdate(r.db.NewUpdate().Model(note)).
		Where("author_id = ?", authorID).
		Where("note_id = ?", noteID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrNoteNotFound
	}

	return note, nil
}

func NewUpdateNoteRepository(db bun.IDB) UpdateNoteRepository {
	return &updateNoteRepositoryImpl{
		db: db,
	}
}
