package dao_test

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/stretchr/testify/require"
	"testing"
)

var updateMessageFixtures = []*entities.Message{
	{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000001",
		Content:    "Lorem Ipsum",
		TargetName: "foo",
	},
}

func TestUpdateMessage(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		teamID    string
		messageID string
		data      *dao.UpdateMessageData

		expect    *entities.Message
		expectErr error
	}{
		{
			name:      "UpdateMessage",
			teamID:    "00000000-0000-0000-0000-000000000001",
			messageID: "00000000-0000-0000-0000-000000000001",
			data: &dao.UpdateMessageData{
				MessageContent: "Lorem ipsum dolor sit amet",
				TargetName:     "foo",
			},
			expect: &entities.Message{
				TeamID:     "00000000-0000-0000-0000-000000000001",
				MessageID:  "00000000-0000-0000-0000-000000000001",
				Content:    "Lorem ipsum dolor sit amet",
				TargetName: "foo",
			},
		},
		{
			name:      "Error/MessageNotFound",
			teamID:    "00000000-0000-0000-0000-000000000001",
			messageID: "00000000-0000-0000-0000-000000000002",
			data: &dao.UpdateMessageData{
				MessageContent: "Lorem Ipsum",
				TargetName:     "foo",
			},
			expectErr: dao.ErrMessageNotFound,
		},
	}

	stx := BeginTX(db, updateMessageFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewUpdateMessageRepository(tx)
			note, err := repo.UpdateMessage(context.TODO(), tt.teamID, tt.messageID, tt.data)

			if note != nil {
				// Since ID is random, nullify if for comparison.
				note.ID = nil
			}

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, note)
		})
	}
}
