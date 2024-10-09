package services_test

import (
	"context"
	daomocks "github.com/in-rich/uservice-search/pkg/dao/mocks"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearchNotes(t *testing.T) {
	testData := []struct {
		name string

		note *models.SearchNotes

		shouldCallSearch    bool
		searchNotesResponse []*entities.Note
		searchNotesError    error

		expect    []*models.Note
		expectErr error
	}{
		{
			name: "Success/Cat",
			note: &models.SearchNotes{
				AuthorID: "00000000-0000-0000-0000-000000000001",
				Limit:    1000,
				Offset:   10,
				RawQuery: "cat",
			},
			searchNotesResponse: []*entities.Note{
				{
					AuthorID:   "00000000-0000-0000-0000-000000000001",
					NoteID:     "00000000-0000-0000-0000-000000000002",
					TargetName: "foo bar",
					Content:    "content",
				},
			},
			expect: []*models.Note{
				{
					AuthorID: "00000000-0000-0000-0000-000000000001",
					NoteID:   "00000000-0000-0000-0000-000000000002",
				},
			},
			shouldCallSearch: true,
		},
		{
			name: "Error/NoteNotFound",
			note: &models.SearchNotes{
				AuthorID: "00000000-0000-0000-0000-000000000001",
				Limit:    1000,
				Offset:   0,
				RawQuery: "cat",
			},
			searchNotesError: FooErr,
			expectErr:        FooErr,
			shouldCallSearch: true,
		},
		{
			name: "Error/Invalid",
			note: &models.SearchNotes{
				AuthorID: "00000000-0000-0000-0000-000000000001",
				Limit:    -12,
				Offset:   0,
				RawQuery: "cat",
			},
			searchNotesError: services.ErrInvalidNoteSearch,
			expectErr:        services.ErrInvalidNoteSearch,
			shouldCallSearch: false,
		},
		{
			name: "Error/FooErr",
			note: &models.SearchNotes{
				AuthorID: "00000000-0000-0000-0000-000000000001",
				Limit:    1000,
				Offset:   0,
				RawQuery: "cat",
			},
			searchNotesError: FooErr,
			expectErr:        FooErr,
			shouldCallSearch: true,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			searchNotes := daomocks.NewMockSearchNotesRepository(t)

			if tt.shouldCallSearch {
				searchNotes.On(
					"SearchNotes",
					context.TODO(),
					tt.note.AuthorID,
					tt.note.RawQuery,
					tt.note.Limit,
					tt.note.Offset,
				).Return(tt.searchNotesResponse, tt.searchNotesError)
			}

			service := services.NewSearchNotesService(searchNotes)

			note, err := service.Exec(context.TODO(), tt.note)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)

			searchNotes.AssertExpectations(t)
		})
	}
}
