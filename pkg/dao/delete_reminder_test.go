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

var deleteReminderFixtures = []*entities.Reminder{
	{
		AuthorID:   "authorID-1",
		ReminderID: "00000000-0000-0000-0000-000000000001",
		ExpiredAt:  lo.ToPtr(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)),
	},
}

func TestDeleteReminder(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		authorID   string
		reminderID string

		expectErr error
	}{
		{
			name:       "DeleteReminder",
			authorID:   "authorID-1",
			reminderID: "00000000-0000-0000-0000-000000000001",
		},
		{
			// Still a success because this method if forgiving.
			name:       "ReminderNotFound",
			authorID:   "authorID-1",
			reminderID: "00000000-0000-0000-0000-000000000002",
		},
	}

	stx := BeginTX(db, deleteReminderFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewDeleteReminderRepository(tx)
			err := repo.DeleteReminder(context.TODO(), tt.authorID, tt.reminderID)

			require.ErrorIs(t, err, tt.expectErr)
		})
	}
}
