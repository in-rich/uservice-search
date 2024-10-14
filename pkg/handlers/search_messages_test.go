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

func TestSearchMessage(t *testing.T) {
	testData := []struct {
		name string

		in *search_pb.SearchMessagesRequest

		searchResponse []*models.Message
		searchErr      error

		expect     *search_pb.SearchMessagesResponse
		expectCode codes.Code
	}{
		{
			name: "SearchMessages/Success",
			in: &search_pb.SearchMessagesRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				Search: "cat",
				Limit:  100,
				Offset: 0,
			},
			searchResponse: []*models.Message{
				{
					TeamID:     "00000000-0000-0000-0000-000000000001",
					MessageID:  "00000000-0000-0000-0000-000000000001",
					Content:    "big cat",
					TargetName: "foo bar",
				},
			},
			expect: &search_pb.SearchMessagesResponse{
				Messages: []*search_pb.Message{
					{
						TeamId:    "00000000-0000-0000-0000-000000000001",
						MessageId: "00000000-0000-0000-0000-000000000001",
					},
				},
			},
		},
		{
			name: "SearchMessages/InvalidArgument",
			in: &search_pb.SearchMessagesRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				Search: "cat",
				Limit:  -4,
				Offset: 0,
			},
			searchErr:  services.ErrInvalidMessageSearch,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "SearchMessages/MessageNotFound",
			in: &search_pb.SearchMessagesRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				Search: "cat",
				Limit:  -4,
				Offset: 0,
			},
			searchErr:  errors.New("Message not found"),
			expectCode: codes.Internal,
		},
		{
			name: "Internal",
			in: &search_pb.SearchMessagesRequest{
				TeamId: "00000000-0000-0000-0000-000000000001",
				Search: "cat",
				Limit:  100,
				Offset: 0,
			},
			searchErr:  errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockSearchMessagesService(t)
			service.On("Exec", context.TODO(), &models.SearchMessages{
				TeamID:   tt.in.GetTeamId(),
				Limit:    int(tt.in.GetLimit()),
				Offset:   int(tt.in.GetOffset()),
				RawQuery: tt.in.GetSearch(),
			}).Return(tt.searchResponse, tt.searchErr)

			handler := handlers.NewSearchMessagesHandler(service)

			resp, err := handler.SearchMessages(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
