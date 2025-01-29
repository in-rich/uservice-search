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
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestUpsertReminder(t *testing.T) {
	testData := []struct {
		name string

		in *search_pb.UpsertReminderRequest

		upsertResponse *models.Reminder
		upsertErr      error

		expect     *search_pb.UpsertReminderResponse
		expectCode codes.Code
	}{
		{
			name: "UpsertReminder",
			in: &search_pb.UpsertReminderRequest{
				AuthorId:         "authorID-1",
				ReminderId:       "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
				ExpiredAt:        timestamppb.New(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			upsertResponse: &models.Reminder{
				AuthorID:   "authorID-1",
				ReminderID: "00000000-0000-0000-0000-000000000001",
				Content:    "content",
				TargetName: "foo bar",
			},
			expect: &search_pb.UpsertReminderResponse{
				Reminder: &search_pb.Reminder{
					AuthorId:   "authorID-1",
					ReminderId: "00000000-0000-0000-0000-000000000001",
				},
			},
		},
		{
			name: "DeleteUser",
			in: &search_pb.UpsertReminderRequest{
				AuthorId:         "authorID-1",
				ReminderId:       "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
				ExpiredAt:        timestamppb.New(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &search_pb.UpsertReminderResponse{},
		},
		{
			name: "InvalidArgument",
			in: &search_pb.UpsertReminderRequest{
				AuthorId:         "authorID-1",
				ReminderId:       "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
				ExpiredAt:        timestamppb.New(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			upsertErr:  services.ErrInvalidReminderUpdate,
			expectCode: codes.InvalidArgument,
		},
		{
			name: "Internal",
			in: &search_pb.UpsertReminderRequest{
				AuthorId:         "authorID-1",
				ReminderId:       "00000000-0000-0000-0000-000000000001",
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
			service := servicesmocks.NewMockUpsertReminderService(t)
			service.On("Exec", context.TODO(), &models.UpsertReminder{
				AuthorID:         tt.in.GetAuthorId(),
				ReminderID:       tt.in.GetReminderId(),
				Content:          tt.in.GetContent(),
				TargetName:       tt.in.GetTargetName(),
				PublicIdentifier: tt.in.GetPublicIdentifier(),
				ExpiredAt:        lo.ToPtr(tt.in.GetExpiredAt().AsTime()),
			}).Return(tt.upsertResponse, tt.upsertErr)

			handler := handlers.NewUpsertReminderHandler(service)

			resp, err := handler.UpsertReminder(context.TODO(), tt.in)

			require.Equal(t, tt.expect, resp)
			RequireGRPCCodesEqual(t, err, tt.expectCode)

			service.AssertExpectations(t)
		})
	}
}
