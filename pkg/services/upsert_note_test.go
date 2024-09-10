package services_test

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/dao"
	daomocks "github.com/in-rich/uservice-search/pkg/dao/mocks"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpsertNote(t *testing.T) {
	testData := []struct {
		name string

		note *models.UpsertNote

		shouldCallCountUpdates bool
		countUpdatesResponse   int
		countUpdatesError      error

		shouldCallDeleteNote bool
		deleteNoteError      error

		shouldCallCreateNote bool
		createNoteResponse   *entities.Note
		createNoteError      error

		shouldCallUpdateNote bool
		updateNoteResponse   *entities.Note
		updateNoteError      error

		expect    *models.Note
		expectErr error
	}{
		{
			name: "UpdateNote",
			note: &models.UpsertNote{
				AuthorID:         "00000000-0000-0000-0000-000000000001",
				NoteID:           "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallCreateNote: true,
			createNoteError:      dao.ErrNoteAlreadyExists,
			shouldCallUpdateNote: true,
			updateNoteResponse: &entities.Note{
				AuthorID:   "00000000-0000-0000-0000-000000000001",
				NoteID:     "00000000-0000-0000-0000-000000000001",
				TargetName: "foo bar",
				Content:    "content",
			},
			expect: &models.Note{
				AuthorID:   "00000000-0000-0000-0000-000000000001",
				NoteID:     "00000000-0000-0000-0000-000000000001",
				TargetName: "foo bar",
				Content:    "content",
			},
		},
		{
			name: "CreateNote",
			note: &models.UpsertNote{
				AuthorID:         "00000000-0000-0000-0000-000000000001",
				NoteID:           "00000000-0000-0000-0000-000000000002",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallCreateNote: true,
			createNoteResponse: &entities.Note{
				AuthorID:   "00000000-0000-0000-0000-000000000001",
				NoteID:     "00000000-0000-0000-0000-000000000002",
				TargetName: "foo bar",
				Content:    "content",
			},
			expect: &models.Note{
				AuthorID:   "00000000-0000-0000-0000-000000000001",
				NoteID:     "00000000-0000-0000-0000-000000000002",
				TargetName: "foo bar",
				Content:    "content",
			},
		},
		{
			name: "DeleteNote",
			note: &models.UpsertNote{
				AuthorID:         "00000000-0000-0000-0000-000000000001",
				NoteID:           "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallDeleteNote: true,
		},
		{
			name: "UpdateNoteError",
			note: &models.UpsertNote{
				AuthorID:         "00000000-0000-0000-0000-000000000001",
				NoteID:           "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallCreateNote: true,
			createNoteError:      dao.ErrNoteAlreadyExists,
			shouldCallUpdateNote: true,
			updateNoteError:      FooErr,
			expectErr:            FooErr,
		},
		{
			name: "CreateNoteError",
			note: &models.UpsertNote{
				AuthorID:         "00000000-0000-0000-0000-000000000001",
				NoteID:           "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallCreateNote: true,
			createNoteError:      FooErr,
			expectErr:            FooErr,
		},
		{
			name: "DeleteNoteError",
			note: &models.UpsertNote{
				AuthorID:         "00000000-0000-0000-0000-000000000001",
				NoteID:           "00000000-0000-0000-0000-000000000001",
				Content:          "",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallDeleteNote: true,
			deleteNoteError:      FooErr,
			expectErr:            FooErr,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			deleteNote := daomocks.NewMockDeleteNoteRepository(t)
			createNote := daomocks.NewMockCreateNoteRepository(t)
			updateNote := daomocks.NewMockUpdateNoteRepository(t)

			if tt.shouldCallDeleteNote {
				deleteNote.
					On("DeleteNote", context.TODO(), tt.note.AuthorID, tt.note.NoteID).
					Return(tt.deleteNoteError)
			}

			if tt.shouldCallCreateNote {
				createNote.
					On(
						"CreateNote",
						context.TODO(),
						tt.note.AuthorID,
						tt.note.NoteID,
						&dao.CreateNoteData{
							NoteContent: tt.note.Content,
							TargetName:  tt.note.TargetName + " " + tt.note.PublicIdentifier,
						},
					).
					Return(tt.createNoteResponse, tt.createNoteError)
			}

			if tt.shouldCallUpdateNote {
				updateNote.
					On(
						"UpdateNote",
						context.TODO(),
						tt.note.AuthorID,
						tt.note.NoteID,
						&dao.UpdateNoteData{
							NoteContent: tt.note.Content,
							TargetName:  tt.note.TargetName + " " + tt.note.PublicIdentifier,
						},
					).
					Return(tt.updateNoteResponse, tt.updateNoteError)
			}

			service := services.NewUpsertNoteService(updateNote, createNote, deleteNote)

			note, err := service.Exec(context.TODO(), tt.note)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)

			deleteNote.AssertExpectations(t)
			createNote.AssertExpectations(t)
			updateNote.AssertExpectations(t)
		})
	}
}
