package services_test

import (
	"context"
	daomocks "github.com/in-rich/uservice-search/pkg/dao/mocks"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearchReminders(t *testing.T) {
	testData := []struct {
		name string

		reminder *models.SearchReminders

		shouldCallSearch        bool
		searchRemindersResponse []*entities.Reminder
		searchRemindersError    error

		expect    []*models.Reminder
		expectErr error
	}{
		{
			name: "Success/Cat",
			reminder: &models.SearchReminders{
				AuthorID: "00000000-0000-0000-0000-000000000001",
				Limit:    1000,
				Offset:   10,
				RawQuery: "cat",
			},
			searchRemindersResponse: []*entities.Reminder{
				{
					AuthorID:   "00000000-0000-0000-0000-000000000001",
					ReminderID: "00000000-0000-0000-0000-000000000002",
					TargetName: "foo bar",
					Content:    "content",
				},
			},
			expect: []*models.Reminder{
				{
					AuthorID:   "00000000-0000-0000-0000-000000000001",
					ReminderID: "00000000-0000-0000-0000-000000000002",
				},
			},
			shouldCallSearch: true,
		},
		{
			name: "Error/ReminderNotFound",
			reminder: &models.SearchReminders{
				AuthorID: "00000000-0000-0000-0000-000000000001",
				Limit:    1000,
				Offset:   0,
				RawQuery: "cat",
			},
			searchRemindersError: FooErr,
			expectErr:            FooErr,
			shouldCallSearch:     true,
		},
		{
			name: "Error/Invalid",
			reminder: &models.SearchReminders{
				AuthorID: "00000000-0000-0000-0000-000000000001",
				Limit:    -12,
				Offset:   0,
				RawQuery: "cat",
			},
			searchRemindersError: services.ErrInvalidReminderSearch,
			expectErr:            services.ErrInvalidReminderSearch,
			shouldCallSearch:     false,
		},
		{
			name: "Error/FooErr",
			reminder: &models.SearchReminders{
				AuthorID: "00000000-0000-0000-0000-000000000001",
				Limit:    1000,
				Offset:   0,
				RawQuery: "cat",
			},
			searchRemindersError: FooErr,
			expectErr:            FooErr,
			shouldCallSearch:     true,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			searchReminders := daomocks.NewMockSearchRemindersRepository(t)

			if tt.shouldCallSearch {
				searchReminders.On(
					"SearchReminders",
					context.TODO(),
					tt.reminder.AuthorID,
					tt.reminder.RawQuery,
					tt.reminder.Limit,
					tt.reminder.Offset,
				).Return(tt.searchRemindersResponse, tt.searchRemindersError)
			}

			service := services.NewSearchRemindersService(searchReminders)

			reminder, err := service.Exec(context.TODO(), tt.reminder)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, reminder)

			searchReminders.AssertExpectations(t)
		})
	}
}
