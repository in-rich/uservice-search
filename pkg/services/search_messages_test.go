package services_test

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/dao"
	daomocks "github.com/in-rich/uservice-search/pkg/dao/mocks"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearchMessages(t *testing.T) {
	testData := []struct {
		name string

		Message *models.SearchMessages

		shouldCallSearch       bool
		searchMessagesResponse []*entities.Message
		searchMessagesError    error

		expect    []*models.Message
		expectErr error
	}{
		{
			name: "Success/Cat",
			Message: &models.SearchMessages{
				UserID:           "00000000-0000-0000-0000-000000000001",
				TeamID:           "",
				Limit:            1000,
				Offset:           10,
				RawQuery:         "cat",
				OneMessageByTeam: true,
			},
			searchMessagesResponse: []*entities.Message{
				{
					TeamID:     "00000000-0000-0000-0000-000000000001",
					MessageID:  "00000000-0000-0000-0000-000000000002",
					TargetName: "foo bar",
					Content:    "content",
				},
			},
			expect: []*models.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000002",
				},
			},
			shouldCallSearch: true,
		},
		{
			name: "Error/Invalid",
			Message: &models.SearchMessages{
				UserID:           "00000000-0000-0000-0000-000000000001",
				TeamID:           "",
				Limit:            -12,
				Offset:           0,
				RawQuery:         "cat",
				OneMessageByTeam: true,
			},
			searchMessagesError: services.ErrInvalidMessageSearch,
			expectErr:           services.ErrInvalidMessageSearch,
			shouldCallSearch:    false,
		},
		{
			name: "Error/FooErr",
			Message: &models.SearchMessages{
				UserID:           "",
				TeamID:           "00000000-0000-0000-0000-000000000001",
				Limit:            1000,
				Offset:           0,
				RawQuery:         "cat",
				OneMessageByTeam: true,
			},
			searchMessagesError: FooErr,
			expectErr:           FooErr,
			shouldCallSearch:    true,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			searchMessages := daomocks.NewMockSearchMessagesRepository(t)

			if tt.shouldCallSearch {
				searchMessages.On(
					"SearchMessages",
					context.TODO(),
					&dao.SearchMessageData{
						UserID:           tt.Message.UserID,
						TeamID:           tt.Message.TeamID,
						RawQuery:         tt.Message.RawQuery,
						Limit:            tt.Message.Limit,
						Offset:           tt.Message.Offset,
						OneMessageByTeam: tt.Message.OneMessageByTeam,
					},
				).Return(tt.searchMessagesResponse, tt.searchMessagesError)
			}

			service := services.NewSearchMessagesService(searchMessages)

			Message, err := service.Exec(context.TODO(), tt.Message)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, Message)

			searchMessages.AssertExpectations(t)
		})
	}
}
