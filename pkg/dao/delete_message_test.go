package dao_test

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/stretchr/testify/require"
	"testing"
)

var deleteMessageFixtures = []*entities.Message{
	{
		TeamID:    "00000000-0000-0000-0000-000000000001",
		MessageID: "00000000-0000-0000-0000-000000000001",
	},
}

func TestDeleteMessage(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		teamID    string
		messageID string

		expectErr error
	}{
		{
			name:      "DeleteMessage",
			teamID:    "00000000-0000-0000-0000-000000000001",
			messageID: "00000000-0000-0000-0000-000000000001",
		},
		{
			// Still a success because this method if forgiving.
			name:      "MessageNotFound",
			teamID:    "00000000-0000-0000-0000-000000000001",
			messageID: "00000000-0000-0000-0000-000000000002",
		},
	}

	stx := BeginTX(db, deleteMessageFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewDeleteMessageRepository(tx)
			err := repo.DeleteMessage(context.TODO(), tt.teamID, tt.messageID)

			require.ErrorIs(t, err, tt.expectErr)
		})
	}
}
