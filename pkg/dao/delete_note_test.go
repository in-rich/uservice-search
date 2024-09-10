package dao_test

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/stretchr/testify/require"
	"testing"
)

var deleteNoteFixtures = []*entities.Note{
	{
		AuthorID: "00000000-0000-0000-0000-000000000001",
		NoteID:   "00000000-0000-0000-0000-000000000001",
	},
}

func TestDeleteNote(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		authorID string
		noteID   string

		expectErr error
	}{
		{
			name:     "DeleteNote",
			authorID: "00000000-0000-0000-0000-000000000001",
			noteID:   "00000000-0000-0000-0000-000000000001",
		},
		{
			// Still a success because this method if forgiving.
			name:     "NoteNotFound",
			authorID: "00000000-0000-0000-0000-000000000001",
			noteID:   "00000000-0000-0000-0000-000000000002",
		},
	}

	stx := BeginTX(db, deleteNoteFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewDeleteNoteRepository(tx)
			err := repo.DeleteNote(context.TODO(), tt.authorID, tt.noteID)

			require.ErrorIs(t, err, tt.expectErr)
		})
	}
}
