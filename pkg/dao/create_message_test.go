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

var createMessageFixtures = []*entities.Message{
	{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000002",
		Content:    "Lorem Ipsum",
		TargetName: "foo",
	},
}

func TestCreateMessage(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		teamID    string
		messageID string
		data      *dao.CreateMessageData

		expect    *entities.Message
		expectErr error
	}{
		{
			name:      "Success",
			teamID:    "00000000-0000-0000-0000-000000000001",
			messageID: "00000000-0000-0000-0000-000000000001",
			data: &dao.CreateMessageData{
				MessageContent: "Lorem Ipsum",
				TargetName:     "foo",
				UpdatedAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expect: &entities.Message{
				TeamID:     "00000000-0000-0000-0000-000000000001",
				MessageID:  "00000000-0000-0000-0000-000000000001",
				Content:    "Lorem Ipsum",
				TargetName: "foo",
				UpdatedAt:  lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name:      "Error/MessageAlreadyExists",
			teamID:    "00000000-0000-0000-0000-000000000001",
			messageID: "00000000-0000-0000-0000-000000000002",
			data: &dao.CreateMessageData{
				MessageContent: "Lorem Ipsum",
				TargetName:     "foo",
				UpdatedAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectErr: dao.ErrMessageAlreadyExists,
		},
		{
			name:      "Success/SameMessageDifferentTeam",
			teamID:    "00000000-0000-0000-0000-000000000002",
			messageID: "00000000-0000-0000-0000-000000000001",
			data: &dao.CreateMessageData{
				MessageContent: "Lorem Ipsum",
				TargetName:     "foo",
				UpdatedAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expect: &entities.Message{
				TeamID:     "00000000-0000-0000-0000-000000000002",
				MessageID:  "00000000-0000-0000-0000-000000000001",
				Content:    "Lorem Ipsum",
				TargetName: "foo",
				UpdatedAt:  lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}

	stx := BeginTX(db, createMessageFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewCreateMessageRepository(tx)
			message, err := repo.CreateMessage(context.TODO(), tt.teamID, tt.messageID, tt.data)

			if message != nil {
				// Since ID is random, nullify if for comparison.
				message.ID = nil
			}

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, message)
		})
	}
}
