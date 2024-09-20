package handlers_test

import (
	"context"
	"errors"
	search_pb "github.com/in-rich/proto/proto-go/search"
	"github.com/in-rich/uservice-search/pkg/handlers"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	servicesmocks "github.com/in-rich/uservice-search/pkg/services/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestSearchNote(t *testing.T) {
	testData := []struct {
		name string

		in *search_pb.SearchNotesRequest

		searchResponse []*models.Note
		searchErr      error

		expect     *search_pb.SearchNotesResponse
		expectCode codes.Code
	}{
		{
			name: "SearchNotes/Success",
			in: &search_pb.SearchNotesRequest{
				AuthorId: "00000000-0000-0000-0000-000000000001",
				Search:   "cat",
				Limit:    100,
				Offset:   0,
			},
			searchResponse: []*models.Note{
				{
					AuthorID:   "00000000-0000-0000-0000-000000000001",
					NoteID:     "00000000-0000-0000-0000-000000000001",
					Content:    "big cat",
					TargetName: "foo bar",
				},
			},
			expect: &search_pb.SearchNotesResponse{
				Notes: []*search_pb.Note{
					{
						AuthorId: "00000000-0000-0000-0000-000000000001",
						NoteId:   "00000000-0000-0000-0000-000000000001",
					},
				},
			},
		},
		{
			name: "SearchNotes/InvalidArgument",
			in: &search_pb.SearchNotesRequest{
				AuthorId: "00000000-0000-0000-0000-000000000001",
				Search:   "cat",
				Limit:    -4,
				Offset:   0,
			},
			searchErr:  services.ErrInvalidNoteSearch,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "SearchNotes/NoteNotFound",
			in: &search_pb.SearchNotesRequest{
				AuthorId: "00000000-0000-0000-0000-000000000001",
				Search:   "cat",
				Limit:    -4,
				Offset:   0,
			},
			searchErr:  errors.New("note not found"),
			expectCode: codes.Internal,
		},
		{
			name: "Internal",
			in: &search_pb.SearchNotesRequest{
				AuthorId: "00000000-0000-0000-0000-000000000001",
				Search:   "cat",
				Limit:    100,
				Offset:   0,
			},
			searchErr:  errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockSearchNotesService(t)
			service.On("Exec", context.TODO(), &models.SearchNotes{
				AuthorID: tt.in.GetAuthorId(),
				Limit:    int(tt.in.GetLimit()),
				Offset:   int(tt.in.GetOffset()),
				RawQuery: tt.in.GetSearch(),
			}).Return(tt.searchResponse, tt.searchErr)

			handler := handlers.NewSearchNotesHandler(service)

			resp, err := handler.SearchNotes(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
