package handlers

import (
	"context"
	"errors"
	search_pb "github.com/in-rich/proto/proto-go/search"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchNoteHandler struct {
	search_pb.SearchNotesServer
	service services.SearchNotesService
}

func (h *SearchNoteHandler) SearchNotes(ctx context.Context, in *search_pb.SearchNotesRequest) (*search_pb.SearchNotesResponse, error) {
	searchNotes, err := h.service.Exec(ctx, &models.SearchNotes{
		AuthorID: in.AuthorId,
		Limit:    int(in.Limit),
		Offset:   int(in.Offset),
		RawQuery: in.Search,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidNoteSearch) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid note search: %v", err)
		}

		if errors.Is(err, services.ErrNoteNotFound) {
			return nil, status.Errorf(codes.NotFound, "note not found: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to upsert note: %v", err)
	}

	var notes []*search_pb.Note
	for _, tt := range searchNotes {
		notes = append(notes, &search_pb.Note{
			NoteId:   tt.NoteID,
			AuthorId: tt.AuthorID,
		})
	}

	return &search_pb.SearchNotesResponse{
		Notes: notes,
	}, nil
}

func NewSearchNotesHandler(service services.SearchNotesService) *SearchNoteHandler {
	return &SearchNoteHandler{
		service: service,
	}
}
