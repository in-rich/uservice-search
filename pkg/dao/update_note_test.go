package dao_test

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/stretchr/testify/require"
	"testing"
)

var updateNoteFixtures = []*entities.Note{
	{
		AuthorID:   "authorID-1",
		NoteID:     "00000000-0000-0000-0000-000000000001",
		Content:    "Lorem Ipsum",
		TargetName: "foo",
	},
}

func TestUpdateNote(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		authorID string
		noteID   string
		data     *dao.UpdateNoteData

		expect    *entities.Note
		expectErr error
	}{
		{
			name:     "UpdateNote",
			authorID: "authorID-1",
			noteID:   "00000000-0000-0000-0000-000000000001",
			data: &dao.UpdateNoteData{
				NoteContent: "Lorem ipsum dolor sit amet",
				TargetName:  "foo",
			},
			expect: &entities.Note{
				AuthorID:   "authorID-1",
				NoteID:     "00000000-0000-0000-0000-000000000001",
				Content:    "Lorem ipsum dolor sit amet",
				TargetName: "foo",
			},
		},
		{
			name:     "Error/NoteNotFound",
			authorID: "authorID-1",
			noteID:   "00000000-0000-0000-0000-000000000002",
			data: &dao.UpdateNoteData{
				NoteContent: "Lorem Ipsum",
				TargetName:  "foo",
			},
			expectErr: dao.ErrNoteNotFound,
		},
	}

	stx := BeginTX(db, updateNoteFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewUpdateNoteRepository(tx)
			note, err := repo.UpdateNote(context.TODO(), tt.authorID, tt.noteID, tt.data)

			if note != nil {
				// Since ID is random, nullify if for comparison.
				note.ID = nil
			}

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)
		})
	}
}
