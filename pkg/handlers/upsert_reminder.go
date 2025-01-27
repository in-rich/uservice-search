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

type UpsertReminderHandler struct {
	search_pb.UpsertReminderServer
	service services.UpsertReminderService
}

func (h *UpsertReminderHandler) UpsertReminder(ctx context.Context, in *search_pb.UpsertReminderRequest) (*search_pb.UpsertReminderResponse, error) {
	reminder, err := h.service.Exec(ctx, &models.UpsertReminder{
		AuthorID:         in.GetAuthorId(),
		ReminderID:       in.GetReminderId(),
		Content:          in.GetContent(),
		TargetName:       in.GetTargetName(),
		PublicIdentifier: in.GetPublicIdentifier(),
		ExpiredAt:        lo.ToPtr(in.GetExpiredAt().AsTime()),
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidReminderUpdate) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid reminder update: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to upsert reminder: %v", err)
	}

	if reminder == nil {
		return &search_pb.UpsertReminderResponse{}, nil
	}

	return &search_pb.UpsertReminderResponse{
		Reminder: &search_pb.Reminder{
			AuthorId:   reminder.AuthorID,
			ReminderId: reminder.ReminderID,
		},
	}, nil
}

func NewUpsertReminderHandler(service services.UpsertReminderService) *UpsertReminderHandler {
	return &UpsertReminderHandler{
		service: service,
	}
}
