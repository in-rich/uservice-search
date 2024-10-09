package dao

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type DeleteNoteRepository interface {
	DeleteNote(ctx context.Context, authorID string, noteID string) error
}

type deleteNoteRepositoryImpl struct {
	db bun.IDB
}

func (r *deleteNoteRepositoryImpl) DeleteNote(ctx context.Context, authorID string, noteID string) error {
	_, err := r.db.NewDelete().
		Model(&entities.Note{}).
		Where("author_id = ?", authorID).
		Where("note_id = ?", noteID).
		Exec(ctx)

	return err
}

func NewDeleteNoteRepository(db bun.IDB) DeleteNoteRepository {
	return &deleteNoteRepositoryImpl{
		db: db,
	}
}
