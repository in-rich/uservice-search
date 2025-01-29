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

type SearchReminderHandler struct {
	search_pb.SearchRemindersServer
	service services.SearchRemindersService
}

func (h *SearchReminderHandler) SearchReminders(ctx context.Context, in *search_pb.SearchRemindersRequest) (*search_pb.SearchRemindersResponse, error) {
	remindersModels, err := h.service.Exec(ctx, &models.SearchReminders{
		AuthorID: in.AuthorId,
		Limit:    int(in.Limit),
		Offset:   int(in.Offset),
		RawQuery: in.Search,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidReminderSearch) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid reminder search: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to upsert reminder: %v", err)
	}

	reminders := lo.Map(remindersModels, func(reminder *models.Reminder, _ int) *search_pb.Reminder {
		return &search_pb.Reminder{
			ReminderId: reminder.ReminderID,
			AuthorId:   reminder.AuthorID,
		}
	})
	return &search_pb.SearchRemindersResponse{
		Reminders: reminders,
	}, nil
}

func NewSearchRemindersHandler(service services.SearchRemindersService) *SearchReminderHandler {
	return &SearchReminderHandler{
		service: service,
	}
}
