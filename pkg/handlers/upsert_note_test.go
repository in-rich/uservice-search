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

func TestUpsertNote(t *testing.T) {
	testData := []struct {
		name string

		in *search_pb.UpsertNoteRequest

		upsertResponse *models.Note
		upsertErr      error

		expect     *search_pb.UpsertNoteResponse
		expectCode codes.Code
	}{
		{
			name: "UpsertNote",
			in: &search_pb.UpsertNoteRequest{
				AuthorId:         "00000000-0000-0000-0000-000000000001",
				NoteId:           "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			upsertResponse: &models.Note{
				AuthorID:   "00000000-0000-0000-0000-000000000001",
				NoteID:     "00000000-0000-0000-0000-000000000001",
				Content:    "content",
				TargetName: "foo bar",
			},
			expect: &search_pb.UpsertNoteResponse{
				Note: &search_pb.Note{
					AuthorId: "00000000-0000-0000-0000-000000000001",
					NoteId:   "00000000-0000-0000-0000-000000000001",
				},
			},
		},
		{
			name: "DeleteUser",
			in: &search_pb.UpsertNoteRequest{
				AuthorId:         "00000000-0000-0000-0000-000000000001",
				NoteId:           "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			expect: &search_pb.UpsertNoteResponse{},
		},
		{
			name: "InvalidArgument",
			in: &search_pb.UpsertNoteRequest{
				AuthorId:         "00000000-0000-0000-0000-000000000001",
				NoteId:           "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			upsertErr:  services.ErrInvalidNoteUpdate,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "Internal",
			in: &search_pb.UpsertNoteRequest{
				AuthorId:         "00000000-0000-0000-0000-000000000001",
				NoteId:           "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			upsertErr:  errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockUpsertNoteService(t)
			service.On("Exec", context.TODO(), &models.UpsertNote{
				AuthorID:         tt.in.GetAuthorId(),
				NoteID:           tt.in.GetNoteId(),
				Content:          tt.in.GetContent(),
				TargetName:       tt.in.GetTargetName(),
				PublicIdentifier: tt.in.GetPublicIdentifier(),
			}).Return(tt.upsertResponse, tt.upsertErr)

			handler := handlers.NewUpsertNoteHandler(service)

			resp, err := handler.UpsertNote(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
