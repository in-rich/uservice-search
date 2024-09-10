package dao_test

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/stretchr/testify/require"
	"testing"
)

var createNoteFixtures = []*entities.Note{
	{
		AuthorID:   "00000000-0000-0000-0000-000000000001",
		NoteID:     "00000000-0000-0000-0000-000000000002",
		Content:    "Lorem Ipsum",
		TargetName: "foo",
	},
}

func TestCreateNote(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		authorID string
		noteID   string
		data     *dao.CreateNoteData

		expect    *entities.Note
		expectErr error
	}{
		{
			name:     "Success",
			authorID: "00000000-0000-0000-0000-000000000001",
			noteID:   "00000000-0000-0000-0000-000000000001",
			data: &dao.CreateNoteData{
				NoteContent: "Lorem Ipsum",
				TargetName:  "foo",
			},
			expect: &entities.Note{
				AuthorID:   "00000000-0000-0000-0000-000000000001",
				NoteID:     "00000000-0000-0000-0000-000000000001",
				Content:    "Lorem Ipsum",
				TargetName: "foo",
			},
		},
		{
			name:     "Error/NoteAlreadyExists",
			authorID: "00000000-0000-0000-0000-000000000001",
			noteID:   "00000000-0000-0000-0000-000000000002",
			data: &dao.CreateNoteData{
				NoteContent: "Lorem Ipsum",
				TargetName:  "foo",
			},
			expectErr: dao.ErrNoteAlreadyExists,
		},
		{
			name:     "Success/AuthorAlreadyExists",
			authorID: "00000000-0000-0000-0000-000000000002",
			noteID:   "00000000-0000-0000-0000-000000000001",
			data: &dao.CreateNoteData{
				NoteContent: "Lorem Ipsum",
				TargetName:  "foo",
			},
			expect: &entities.Note{
				AuthorID:   "00000000-0000-0000-0000-000000000002",
				NoteID:     "00000000-0000-0000-0000-000000000001",
				Content:    "Lorem Ipsum",
				TargetName: "foo",
			},
		},
	}

	stx := BeginTX(db, createNoteFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewCreateNoteRepository(tx)
			note, err := repo.CreateNote(context.TODO(), tt.authorID, tt.noteID, tt.data)

			if note != nil {
				// Since ID is random, nullify if for comparison.
				note.ID = nil
			}

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)
		})
	}
}
