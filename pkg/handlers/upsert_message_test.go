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

func TestUpsertMessage(t *testing.T) {
	testData := []struct {
		name string

		in *search_pb.UpsertMessageRequest

		upsertResponse *models.Message
		upsertErr      error

		expect     *search_pb.UpsertMessageResponse
		expectCode codes.Code
	}{
		{
			name: "UpsertMessage",
			in: &search_pb.UpsertMessageRequest{
				TeamId:           "00000000-0000-0000-0000-000000000001",
				MessageId:        "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			upsertResponse: &models.Message{
				TeamID:     "00000000-0000-0000-0000-000000000001",
				MessageID:  "00000000-0000-0000-0000-000000000001",
				Content:    "content",
				TargetName: "foo bar",
			},
			expect: &search_pb.UpsertMessageResponse{
				Message: &search_pb.Message{
					TeamId:    "00000000-0000-0000-0000-000000000001",
					MessageId: "00000000-0000-0000-0000-000000000001",
				},
			},
		},
		{
			name: "DeleteTeam",
			in: &search_pb.UpsertMessageRequest{
				TeamId:           "00000000-0000-0000-0000-000000000001",
				MessageId:        "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			expect: &search_pb.UpsertMessageResponse{},
		},
		{
			name: "InvalidArgument",
			in: &search_pb.UpsertMessageRequest{
				TeamId:           "00000000-0000-0000-0000-000000000001",
				MessageId:        "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			upsertErr:  services.ErrInvalidMessageUpdate,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "Internal",
			in: &search_pb.UpsertMessageRequest{
				TeamId:           "00000000-0000-0000-0000-000000000001",
				MessageId:        "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			upsertErr:  errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockUpsertMessageService(t)
			service.On("Exec", context.TODO(), &models.UpsertMessage{
				MessageID:        tt.in.GetTeamId(),
				TeamID:           tt.in.GetMessageId(),
				Content:          tt.in.GetContent(),
				TargetName:       tt.in.GetTargetName(),
				PublicIdentifier: tt.in.GetPublicIdentifier(),
			}).Return(tt.upsertResponse, tt.upsertErr)

			handler := handlers.NewUpsertMessageHandler(service)

			resp, err := handler.UpsertMessage(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
