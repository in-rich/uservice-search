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

type CreateTeamMetaHandler struct {
	search_pb.CreateTeamMetaServer
	service services.CreateTeamMetaService
}

func (h *CreateTeamMetaHandler) CreateTeamMeta(ctx context.Context, in *search_pb.CreateTeamMetaRequest) (*search_pb.CreateTeamMetaResponse, error) {
	teamMeta, err := h.service.Exec(ctx, &models.TeamMeta{
		TeamID: in.GetTeamId(),
		UserID: in.GetUserId(),
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidTeamMetaCreate) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid team meta update: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to create team meta: %v", err)
	}

	return &search_pb.CreateTeamMetaResponse{
		TeamMeta: &search_pb.TeamMeta{
			TeamId: teamMeta.TeamID,
			UserId: teamMeta.UserID,
		},
	}, nil
}

func NewCreateTeamMetaHandler(service services.CreateTeamMetaService) *CreateTeamMetaHandler {
	return &CreateTeamMetaHandler{
		service: service,
	}
}
