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

type SearchMessagesHandler struct {
	search_pb.SearchMessageServer
	service services.SearchMessagesService
}

func (h *SearchMessagesHandler) SearchMessages(ctx context.Context, in *search_pb.SearchMessagesRequest) (*search_pb.SearchMessagesResponse, error) {
	messagesModels, err := h.service.Exec(ctx, &models.SearchMessages{
		TeamID:   in.TeamId,
		Limit:    int(in.Limit),
		Offset:   int(in.Offset),
		RawQuery: in.Search,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidMessageSearch) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid note search: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to upsert note: %v", err)
	}

	messages := lo.Map(messagesModels, func(message *models.Message, _ int) *search_pb.Message {
		return &search_pb.Message{
			TeamId:    message.TeamID,
			MessageId: message.MessageID,
		}
	})

	return &search_pb.SearchMessagesResponse{
		Messages: messages,
	}, nil
}

func NewSearchMessagesHandler(service services.SearchMessagesService) *SearchMessagesHandler {
	return &SearchMessagesHandler{
		service: service,
	}
}
