package handlers

import (
	"context"
	"errors"
	search_pb "github.com/in-rich/proto/proto-go/search"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchNoteHandler struct {
	search_pb.SearchNotesServer
	service services.SearchNotesService
}

func (h *SearchNoteHandler) SearchNotes(ctx context.Context, in *search_pb.SearchNotesRequest) (*search_pb.SearchNotesResponse, error) {
	notesModels, err := h.service.Exec(ctx, &models.SearchNotes{
		AuthorID: in.AuthorId,
		Limit:    int(in.Limit),
		Offset:   int(in.Offset),
		RawQuery: in.Search,
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidNoteSearch) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid note search: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to upsert note: %v", err)
	}

	notes := lo.Map(notesModels, func(note *models.Note, _ int) *search_pb.Note {
		return &search_pb.Note{
			NoteId:   note.NoteID,
			AuthorId: note.AuthorID,
		}
	})

	return &search_pb.SearchNotesResponse{
		Notes: notes,
	}, nil
}

func NewSearchNotesHandler(service services.SearchNotesService) *SearchNoteHandler {
	return &SearchNoteHandler{
		service: service,
	}
}
