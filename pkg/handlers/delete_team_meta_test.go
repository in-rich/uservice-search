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
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestDeleteTeamMeta(t *testing.T) {
	testData := []struct {
		name string

		in *search_pb.DeleteTeamMetaRequest

		DeleteTeamMetaErr error

		expect     *emptypb.Empty
		expectCode codes.Code
	}{
		{
			name: "DeleteTeamMetaMember",
			in: &search_pb.DeleteTeamMetaRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				UserId: "00000000-0000-0000-0000-000000000001",
			},
			expect: new(emptypb.Empty),
		},
		{
			name: "DeleteTeamMeta",
			in: &search_pb.DeleteTeamMetaRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				UserId: "",
			},
			expect: new(emptypb.Empty),
		},
		{
			name: "InvalidArgument",
			in: &search_pb.DeleteTeamMetaRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				UserId: "00000000-0000-0000-0000-000000000001",
			},
			DeleteTeamMetaErr: services.ErrInvalidTeamMetaDelete,
			expectCode:        codes.InvalidArgument,
		},
		{
			name: "InternalError",
			in: &search_pb.DeleteTeamMetaRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				UserId: "00000000-0000-0000-0000-000000000001",
			},
			DeleteTeamMetaErr: errors.New("internal error"),
			expectCode:        codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockDeleteTeamMetaService(t)
			service.
				On("Exec", context.TODO(), &models.DeleteTeamMeta{
					TeamID: tt.in.GetTeamId(),
					UserID: tt.in.GetUserId(),
				}).
				Return(tt.DeleteTeamMetaErr)
			handler := handlers.NewDeleteTeamMetaHandler(service)

			resp, err := handler.DeleteTeamMeta(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
