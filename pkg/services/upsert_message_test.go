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

func TestUpsertMessage(t *testing.T) {
	testData := []struct {
		name string

		message *models.UpsertMessage

		shouldCallDeleteMessage bool
		deleteMessageError      error

		shouldCallCreateMessage bool
		createMessageResponse   *entities.Message
		createMessageError      error

		shouldCallUpdateMessage bool
		updateMessageResponse   *entities.Message
		updateMessageError      error

		expect    *models.Message
		expectErr error
	}{
		{
			name: "UpdateMessage",
			message: &models.UpsertMessage{
				TeamID:           "00000000-0000-0000-0000-000000000001",
				MessageID:        "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallCreateMessage: true,
			createMessageError:      dao.ErrMessageAlreadyExists,
			shouldCallUpdateMessage: true,
			updateMessageResponse: &entities.Message{
				TeamID:     "00000000-0000-0000-0000-000000000001",
				MessageID:  "00000000-0000-0000-0000-000000000001",
				TargetName: "foo bar",
				Content:    "content",
			},
			expect: &models.Message{
				TeamID:     "00000000-0000-0000-0000-000000000001",
				MessageID:  "00000000-0000-0000-0000-000000000001",
				TargetName: "foo bar",
				Content:    "content",
			},
		},
		{
			name: "CreateMessage",
			message: &models.UpsertMessage{
				TeamID:           "00000000-0000-0000-0000-000000000001",
				MessageID:        "00000000-0000-0000-0000-000000000002",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallCreateMessage: true,
			createMessageResponse: &entities.Message{
				TeamID:     "00000000-0000-0000-0000-000000000001",
				MessageID:  "00000000-0000-0000-0000-000000000002",
				TargetName: "foo bar",
				Content:    "content",
			},
			expect: &models.Message{
				TeamID:     "00000000-0000-0000-0000-000000000001",
				MessageID:  "00000000-0000-0000-0000-000000000002",
				TargetName: "foo bar",
				Content:    "content",
			},
		},
		{
			name: "DeleteMessage",
			message: &models.UpsertMessage{
				TeamID:           "00000000-0000-0000-0000-000000000001",
				MessageID:        "00000000-0000-0000-0000-000000000001",
				Content:          "",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallDeleteMessage: true,
		},
		{
			name: "UpdateMessageError",
			message: &models.UpsertMessage{
				TeamID:           "00000000-0000-0000-0000-000000000001",
				MessageID:        "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallCreateMessage: true,
			createMessageError:      dao.ErrMessageAlreadyExists,
			shouldCallUpdateMessage: true,
			updateMessageError:      FooErr,
			expectErr:               FooErr,
		},
		{
			name: "CreateMessageError",
			message: &models.UpsertMessage{
				TeamID:           "00000000-0000-0000-0000-000000000001",
				MessageID:        "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallCreateMessage: true,
			createMessageError:      FooErr,
			expectErr:               FooErr,
		},
		{
			name: "DeleteMessageError",
			message: &models.UpsertMessage{
				TeamID:           "00000000-0000-0000-0000-000000000001",
				MessageID:        "00000000-0000-0000-0000-000000000001",
				Content:          "",
				TargetName:       "foo",
				PublicIdentifier: "bar",
			},
			shouldCallDeleteMessage: true,
			deleteMessageError:      FooErr,
			expectErr:               FooErr,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			deleteMessage := daomocks.NewMockDeleteMessageRepository(t)
			createMessage := daomocks.NewMockCreateMessageRepository(t)
			updateMessage := daomocks.NewMockUpdateMessageRepository(t)

			if tt.shouldCallDeleteMessage {
				deleteMessage.
					On("DeleteMessage", context.TODO(), tt.message.TeamID, tt.message.MessageID).
					Return(tt.deleteMessageError)
			}

			if tt.shouldCallCreateMessage {
				createMessage.
					On(
						"CreateMessage",
						context.TODO(),
						tt.message.TeamID,
						tt.message.MessageID,
						&dao.CreateMessageData{
							MessageContent: tt.message.Content,
							TargetName:     tt.message.TargetName + " " + tt.message.PublicIdentifier,
						},
					).
					Return(tt.createMessageResponse, tt.createMessageError)
			}

			if tt.shouldCallUpdateMessage {
				updateMessage.
					On(
						"UpdateMessage",
						context.TODO(),
						tt.message.TeamID,
						tt.message.MessageID,
						&dao.UpdateMessageData{
							MessageContent: tt.message.Content,
							TargetName:     tt.message.TargetName + " " + tt.message.PublicIdentifier,
						},
					).
					Return(tt.updateMessageResponse, tt.updateMessageError)
			}

			service := services.NewUpsertMessageService(updateMessage, createMessage, deleteMessage)

			message, err := service.Exec(context.TODO(), tt.message)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, message)

			deleteMessage.AssertExpectations(t)
			createMessage.AssertExpectations(t)
			updateMessage.AssertExpectations(t)
		})
	}
}
