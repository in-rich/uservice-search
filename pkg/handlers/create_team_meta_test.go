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

func TestCreateTeamMeta(t *testing.T) {
	testData := []struct {
		name string

		in *search_pb.CreateTeamMetaRequest

		CreateTeamMetaResponse *models.TeamMeta
		CreateTeamMetaErr      error

		expect     *search_pb.CreateTeamMetaResponse
		expectCode codes.Code
	}{
		{
			name: "CreateTeamMeta",
			in: &search_pb.CreateTeamMetaRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				UserId: "00000000-0000-0000-0000-000000000001",
			},
			CreateTeamMetaResponse: &models.TeamMeta{
				TeamID: "00000000-0000-0000-0000-000000000001",
				UserID: "00000000-0000-0000-0000-000000000001",
			},
			expect: &search_pb.CreateTeamMetaResponse{
				TeamMeta: &search_pb.TeamMeta{
					TeamId: "00000000-0000-0000-0000-000000000001",
					UserId: "00000000-0000-0000-0000-000000000001",
				},
			},
		},
		{
			name: "InvalidArgument",
			in: &search_pb.CreateTeamMetaRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				UserId: "00000000-0000-0000-0000-000000000001",
			},
			CreateTeamMetaErr: services.ErrInvalidTeamMetaCreate,
			expectCode:        codes.InvalidArgument,
		},
		{
			name: "InternalError",
			in: &search_pb.CreateTeamMetaRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				UserId: "00000000-0000-0000-0000-000000000001",
			},
			CreateTeamMetaErr: errors.New("internal error"),
			expectCode:        codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockCreateTeamMetaService(t)
			service.
				On("Exec", context.TODO(), &models.TeamMeta{
					TeamID: "00000000-0000-0000-0000-000000000001",
					UserID: "00000000-0000-0000-0000-000000000001",
				}).
				Return(tt.CreateTeamMetaResponse, tt.CreateTeamMetaErr)
			handler := handlers.NewCreateTeamMetaHandler(service)

			resp, err := handler.CreateTeamMeta(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
