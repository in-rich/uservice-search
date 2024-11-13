package handlers

import (
	"context"
	"errors"
	search_pb "github.com/in-rich/proto/proto-go/search"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DeleteTeamMetaHandler struct {
	search_pb.DeleteTeamMetaServer
	service services.DeleteTeamMetaService
}

func (h *DeleteTeamMetaHandler) DeleteTeamMeta(ctx context.Context, in *search_pb.DeleteTeamMetaRequest) (*emptypb.Empty, error) {
	err := h.service.Exec(ctx, &models.DeleteTeamMeta{
		TeamID: in.GetTeamId(),
		UserID: in.GetUserId(),
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidTeamMetaDelete) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid team meta update: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to delete team meta: %v", err)
	}

	return new(emptypb.Empty), nil
}

func NewDeleteTeamMetaHandler(service services.DeleteTeamMetaService) *DeleteTeamMetaHandler {
	return &DeleteTeamMetaHandler{
		service: service,
	}
}
