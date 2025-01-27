package dao_test

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var updateReminderFixtures = []*entities.Reminder{
	{
		AuthorID:   "authorID-1",
		ReminderID: "00000000-0000-0000-0000-000000000001",
		Content:    "Lorem Ipsum",
		TargetName: "foo",
		ExpiredAt:  lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
	},
}

func TestUpdateReminder(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		authorID   string
		reminderID string
		data       *dao.UpdateReminderData

		expect    *entities.Reminder
		expectErr error
	}{
		{
			name:       "UpdateReminder",
			authorID:   "authorID-1",
			reminderID: "00000000-0000-0000-0000-000000000001",
			data: &dao.UpdateReminderData{
				ReminderContent: "Lorem ipsum dolor sit amet",
				TargetName:      "foo",
				ExpiredAt:       lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &entities.Reminder{
				AuthorID:   "authorID-1",
				ReminderID: "00000000-0000-0000-0000-000000000001",
				Content:    "Lorem ipsum dolor sit amet",
				TargetName: "foo",
				ExpiredAt:  lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name:       "UpdateReminderDate",
			authorID:   "authorID-1",
			reminderID: "00000000-0000-0000-0000-000000000001",
			data: &dao.UpdateReminderData{
				ReminderContent: "Lorem Ipsum",
				TargetName:      "foo",
				ExpiredAt:       lo.ToPtr(time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)),
			},
			expect: &entities.Reminder{
				AuthorID:   "authorID-1",
				ReminderID: "00000000-0000-0000-0000-000000000001",
				Content:    "Lorem Ipsum",
				TargetName: "foo",
				ExpiredAt:  lo.ToPtr(time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name:       "Error/ReminderNotFound",
			authorID:   "authorID-1",
			reminderID: "00000000-0000-0000-0000-000000000002",
			data: &dao.UpdateReminderData{
				ReminderContent: "Lorem Ipsum",
				TargetName:      "foo",
				ExpiredAt:       lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectErr: dao.ErrReminderNotFound,
		},
	}

	stx := BeginTX(db, updateReminderFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewUpdateReminderRepository(tx)
			reminder, err := repo.UpdateReminder(context.TODO(), tt.authorID, tt.reminderID, tt.data)

			if reminder != nil {
				// Since ID is random, nullify if for comparison.
				reminder.ID = nil
				reminder.UpdatedAt = nil
			}

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, reminder)
		})
	}
}
