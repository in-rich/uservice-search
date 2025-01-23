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
	"time"
)

type UpsertNoteHandler struct {
	search_pb.UpsertNoteServer
	service services.UpsertNoteService
}

func (h *UpsertNoteHandler) UpsertNote(ctx context.Context, in *search_pb.UpsertNoteRequest) (*search_pb.UpsertNoteResponse, error) {
	note, err := h.service.Exec(ctx, &models.UpsertNote{
		AuthorID:         in.GetAuthorId(),
		NoteID:           in.GetNoteId(),
		Content:          in.GetContent(),
		TargetName:       in.GetTargetName(),
		PublicIdentifier: in.GetPublicIdentifier(),
		UpdatedAt: lo.TernaryF[*time.Time](
			in.GetUpdatedAt() == nil,
			func() *time.Time {
				return nil
			},
			func() *time.Time {
				return lo.ToPtr(in.GetUpdatedAt().AsTime())
			},
		),
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidNoteUpdate) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid note update: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to upsert note: %v", err)
	}

	if note == nil {
		return &search_pb.UpsertNoteResponse{}, nil
	}

	return &search_pb.UpsertNoteResponse{
		Note: &search_pb.Note{
			AuthorId: note.AuthorID,
			NoteId:   note.NoteID,
		},
	}, nil
}

func NewUpsertNoteHandler(service services.UpsertNoteService) *UpsertNoteHandler {
	return &UpsertNoteHandler{
		service: service,
	}
}
