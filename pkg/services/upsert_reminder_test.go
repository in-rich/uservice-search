package services_test

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/dao"
	daomocks "github.com/in-rich/uservice-search/pkg/dao/mocks"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUpsertReminder(t *testing.T) {
	testData := []struct {
		name string

		reminder *models.UpsertReminder

		shouldCallCountUpdates bool
		countUpdatesResponse   int
		countUpdatesError      error

		shouldCallDeleteReminder bool
		deleteReminderError      error

		shouldCallCreateReminder bool
		createReminderResponse   *entities.Reminder
		createReminderError      error

		shouldCallUpdateReminder bool
		updateReminderResponse   *entities.Reminder
		updateReminderError      error

		expect    *models.Reminder
		expectErr error
	}{
		{
			name: "UpdateReminder",
			reminder: &models.UpsertReminder{
				AuthorID:         "authorID-1",
				ReminderID:       "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
				ExpiredAt:        lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			shouldCallCreateReminder: true,
			createReminderError:      dao.ErrReminderAlreadyExists,
			shouldCallUpdateReminder: true,
			updateReminderResponse: &entities.Reminder{
				AuthorID:   "authorID-1",
				ReminderID: "00000000-0000-0000-0000-000000000001",
				TargetName: "foo bar",
				Content:    "content",
				ExpiredAt:  lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &models.Reminder{
				AuthorID:   "authorID-1",
				ReminderID: "00000000-0000-0000-0000-000000000001",
				TargetName: "foo bar",
				Content:    "content",
			},
		},
		{
			name: "CreateReminder",
			reminder: &models.UpsertReminder{
				AuthorID:         "authorID-1",
				ReminderID:       "00000000-0000-0000-0000-000000000002",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
				ExpiredAt:        lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			shouldCallCreateReminder: true,
			createReminderResponse: &entities.Reminder{
				AuthorID:   "authorID-1",
				ReminderID: "00000000-0000-0000-0000-000000000002",
				TargetName: "foo bar",
				Content:    "content",
				ExpiredAt:  lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &models.Reminder{
				AuthorID:   "authorID-1",
				ReminderID: "00000000-0000-0000-0000-000000000002",
				TargetName: "foo bar",
				Content:    "content",
			},
		},
		{
			name: "DeleteReminder",
			reminder: &models.UpsertReminder{
				AuthorID:         "authorID-1",
				ReminderID:       "00000000-0000-0000-0000-000000000001",
				Content:          "",
				TargetName:       "foo",
				PublicIdentifier: "bar",
				ExpiredAt:        lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			shouldCallDeleteReminder: true,
		},
		{
			name: "UpdateReminderError",
			reminder: &models.UpsertReminder{
				AuthorID:         "authorID-1",
				ReminderID:       "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
				ExpiredAt:        lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			shouldCallCreateReminder: true,
			createReminderError:      dao.ErrReminderAlreadyExists,
			shouldCallUpdateReminder: true,
			updateReminderError:      FooErr,
			expectErr:                FooErr,
		},
		{
			name: "CreateReminderError",
			reminder: &models.UpsertReminder{
				AuthorID:         "authorID-1",
				ReminderID:       "00000000-0000-0000-0000-000000000001",
				Content:          "content",
				TargetName:       "foo",
				PublicIdentifier: "bar",
				ExpiredAt:        lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			shouldCallCreateReminder: true,
			createReminderError:      FooErr,
			expectErr:                FooErr,
		},
		{
			name: "DeleteReminderError",
			reminder: &models.UpsertReminder{
				AuthorID:         "authorID-1",
				ReminderID:       "00000000-0000-0000-0000-000000000001",
				Content:          "",
				TargetName:       "foo",
				PublicIdentifier: "bar",
				ExpiredAt:        lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			shouldCallDeleteReminder: true,
			deleteReminderError:      FooErr,
			expectErr:                FooErr,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			deleteReminder := daomocks.NewMockDeleteReminderRepository(t)
			createReminder := daomocks.NewMockCreateReminderRepository(t)
			updateReminder := daomocks.NewMockUpdateReminderRepository(t)

			if tt.shouldCallDeleteReminder {
				deleteReminder.
					On("DeleteReminder", context.TODO(), tt.reminder.AuthorID, tt.reminder.ReminderID).
					Return(tt.deleteReminderError)
			}

			if tt.shouldCallCreateReminder {
				createReminder.
					On(
						"CreateReminder",
						context.TODO(),
						tt.reminder.AuthorID,
						tt.reminder.ReminderID,
						&dao.CreateReminderData{
							ReminderContent: tt.reminder.Content,
							TargetName:      tt.reminder.TargetName + " " + tt.reminder.PublicIdentifier,
							ExpiredAt:       tt.reminder.ExpiredAt,
						},
					).
					Return(tt.createReminderResponse, tt.createReminderError)
			}

			if tt.shouldCallUpdateReminder {
				updateReminder.
					On(
						"UpdateReminder",
						context.TODO(),
						tt.reminder.AuthorID,
						tt.reminder.ReminderID,
						&dao.UpdateReminderData{
							ReminderContent: tt.reminder.Content,
							TargetName:      tt.reminder.TargetName + " " + tt.reminder.PublicIdentifier,
							ExpiredAt:       tt.reminder.ExpiredAt,
						},
					).
					Return(tt.updateReminderResponse, tt.updateReminderError)
			}

			service := services.NewUpsertReminderService(updateReminder, createReminder, deleteReminder)

			reminder, err := service.Exec(context.TODO(), tt.reminder)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, reminder)

			deleteReminder.AssertExpectations(t)
			createReminder.AssertExpectations(t)
			updateReminder.AssertExpectations(t)
		})
	}
}
