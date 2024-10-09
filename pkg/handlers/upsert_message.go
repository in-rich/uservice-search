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

type UpsertMessageHandler struct {
	search_pb.UpsertMessageServer
	service services.UpsertMessageService
}

func (h *UpsertMessageHandler) UpsertMessage(ctx context.Context, in *search_pb.UpsertMessageRequest) (*search_pb.UpsertMessageResponse, error) {
	message, err := h.service.Exec(ctx, &models.UpsertMessage{
		TeamID:           in.GetTeamId(),
		MessageID:        in.GetMessageId(),
		Content:          in.GetContent(),
		TargetName:       in.GetTargetName(),
		PublicIdentifier: in.GetPublicIdentifier(),
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidMessageUpdate) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid message update: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to upsert message: %v", err)
	}

	if message == nil {
		return &search_pb.UpsertMessageResponse{}, nil
	}

	return &search_pb.UpsertMessageResponse{
		Message: &search_pb.Message{
			TeamId:    message.TeamID,
			MessageId: message.MessageID,
		},
	}, nil
}

func NewUpsertMessageHandler(service services.UpsertMessageService) *UpsertMessageHandler {
	return &UpsertMessageHandler{
		service: service,
	}
}
