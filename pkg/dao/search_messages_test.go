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

var searchMessageFixtures = []*entities.Message{
	{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000001",
		Content:    "big cat and rat",
		TargetName: "John Smith",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 1, 0, 0, 0, time.UTC)),
	},
	{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000002",
		Content:    "fat rat",
		TargetName: "Xavier Login",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 2, 0, 0, 0, time.UTC)),
	},
	{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000003",
		Content:    "I like train",
		TargetName: "Jane Doe",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 3, 0, 0, 0, time.UTC)),
	},
	{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000004",
		Content:    "blue rat and green cat",
		TargetName: "Alice",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 4, 0, 0, 0, time.UTC)),
	},
	{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000005",
		Content:    "Lorem Ipsum",
		TargetName: "Bob",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 5, 0, 0, 0, time.UTC)),
	},
}

func TestSearchMessage(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		TeamID   string
		rawQuery string

		expect    []*entities.Message
		expectErr error
	}{
		{
			name:     "Success/RatAndCat",
			TeamID:   "00000000-0000-0000-0000-000000000001",
			rawQuery: "cat rat",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000001",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000004",
				},
			},
		},
		{
			name:     "Success/Rat",
			TeamID:   "00000000-0000-0000-0000-000000000001",
			rawQuery: "rat",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000004",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000002",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000001",
				},
			},
		},
		{
			name:     "Success/Empty",
			TeamID:   "00000000-0000-0000-0000-000000000001",
			rawQuery: "",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000005",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000004",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000003",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000002",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000001",
				},
			},
		},
	}

	stx := BeginTX(db, searchMessageFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {

			repo := dao.NewSearchMessagesRepository(stx)
			Message, err := repo.SearchMessages(context.TODO(), tt.TeamID, tt.rawQuery, 1000, 0)

			for _, tt := range Message {
				// Since ID and UpdatedAt are random, nullify them for comparison.
				tt.ID = nil
				tt.UpdatedAt = nil

				// Since Content and TargetName are not relevant at this point we empty them.
				tt.Content = ""
				tt.TargetName = ""
			}

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, Message)
		})
	}
}
