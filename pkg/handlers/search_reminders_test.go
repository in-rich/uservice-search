package handlers_test

import (
	"context"
	"errors"
	search_pb "github.com/in-rich/proto/proto-go/search"
	"github.com/in-rich/uservice-search/pkg/handlers"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	servicesmocks "github.com/in-rich/uservice-search/pkg/services/mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
	"time"
)

func TestSearchReminder(t *testing.T) {
	testData := []struct {
		name string

		in *search_pb.SearchRemindersRequest

		searchResponse []*models.Reminder
		searchErr      error

		expect     *search_pb.SearchRemindersResponse
		expectCode codes.Code
	}{
		{
			name: "SearchReminders/Success",
			in: &search_pb.SearchRemindersRequest{
				AuthorId: "00000000-0000-0000-0000-000000000001",
				Search:   "cat",
				Limit:    100,
				Offset:   0,
			},
			searchResponse: []*models.Reminder{
				{
					AuthorID:   "00000000-0000-0000-0000-000000000001",
					ReminderID: "00000000-0000-0000-0000-000000000001",
					Content:    "big cat",
					TargetName: "foo bar",
					ExpiredAt:  lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			expect: &search_pb.SearchRemindersResponse{
				Reminders: []*search_pb.Reminder{
					{
						AuthorId:   "00000000-0000-0000-0000-000000000001",
						ReminderId: "00000000-0000-0000-0000-000000000001",
					},
				},
			},
		},
		{
			name: "SearchReminders/InvalidArgument",
			in: &search_pb.SearchRemindersRequest{
				AuthorId: "00000000-0000-0000-0000-000000000001",
				Search:   "cat",
				Limit:    -4,
				Offset:   0,
			},
			searchErr:  services.ErrInvalidReminderSearch,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "SearchReminders/ReminderNotFound",
			in: &search_pb.SearchRemindersRequest{
				AuthorId: "00000000-0000-0000-0000-000000000001",
				Search:   "cat",
				Limit:    -4,
				Offset:   0,
			},
			searchErr:  errors.New("reminder not found"),
			expectCode: codes.Internal,
		},
		{
			name: "Internal",
			in: &search_pb.SearchRemindersRequest{
				AuthorId: "00000000-0000-0000-0000-000000000001",
				Search:   "cat",
				Limit:    100,
				Offset:   0,
			},
			searchErr:  errors.New("internal error"),
			expectCode: codes.Internal,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			service := servicesmocks.NewMockSearchRemindersService(t)
			service.On("Exec", context.TODO(), &models.SearchReminders{
				AuthorID: tt.in.GetAuthorId(),
				Limit:    int(tt.in.GetLimit()),
				Offset:   int(tt.in.GetOffset()),
				RawQuery: tt.in.GetSearch(),
			}).Return(tt.searchResponse, tt.searchErr)

			handler := handlers.NewSearchRemindersHandler(service)

			resp, err := handler.SearchReminders(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
