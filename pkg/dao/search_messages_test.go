package dao_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var searchMessageFixtures = []interface{}{
	&entities.Message{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000001",
		Content:    "big cat and rat",
		TargetName: "John Smith",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 1, 0, 0, 0, time.UTC)),
	},
	&entities.Message{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000002",
		Content:    "fat rat",
		TargetName: "Xavier Login",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 2, 0, 0, 0, time.UTC)),
	},
	&entities.Message{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000003",
		Content:    "I like train",
		TargetName: "Jane Doe",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 3, 0, 0, 0, time.UTC)),
	},
	&entities.Message{
		TeamID:     "00000000-0000-0000-0000-000000000002",
		MessageID:  "00000000-0000-0000-0000-000000000004",
		Content:    "blue rat and green cat",
		TargetName: "Alice",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 4, 0, 0, 0, time.UTC)),
	},
	&entities.Message{
		TeamID:     "00000000-0000-0000-0000-000000000002",
		MessageID:  "00000000-0000-0000-0000-000000000005",
		Content:    "Lorem Ipsum",
		TargetName: "Bob",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 5, 0, 0, 0, time.UTC)),
	},
	&entities.Message{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000006",
		Content:    "foo",
		TargetName: "John Smith",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 6, 0, 0, 0, time.UTC)),
	},
	&entities.Message{
		TeamID:     "00000000-0000-0000-0000-000000000001",
		MessageID:  "00000000-0000-0000-0000-000000000007",
		Content:    "hello world",
		TargetName: "John Smith",
		UpdatedAt:  lo.ToPtr(time.Date(2024, time.September, 17, 7, 0, 0, 0, time.UTC)),
	},
	&entities.TeamMeta{
		TeamID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		UserID: "00000000-0000-0000-0000-000000000001",
	},
	&entities.TeamMeta{
		TeamID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
		UserID: "00000000-0000-0000-0000-000000000001",
	},
	&entities.TeamMeta{
		TeamID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		UserID: "00000000-0000-0000-0000-000000000002",
	},
	&entities.TeamMeta{
		TeamID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
		UserID: "00000000-0000-0000-0000-000000000003",
	},
}

func TestSearchMessage(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		UserID           string
		TeamID           string
		rawQuery         string
		OneMessageByTeam bool

		expect    []*entities.Message
		expectErr error
	}{
		{
			name:             "Success/RatAndCatUser1",
			UserID:           "00000000-0000-0000-0000-000000000001",
			TeamID:           "",
			OneMessageByTeam: false,
			rawQuery:         "cat rat",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000001",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000002",
					MessageID: "00000000-0000-0000-0000-000000000004",
				},
			},
		},
		{
			name:             "Success/RatAndCatUser1OneMessage",
			UserID:           "00000000-0000-0000-0000-000000000001",
			TeamID:           "",
			OneMessageByTeam: true,
			rawQuery:         "cat rat",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000001",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000002",
					MessageID: "00000000-0000-0000-0000-000000000004",
				},
			},
		},
		{
			name:             "Success/RatAndCatUser2",
			UserID:           "00000000-0000-0000-0000-000000000002",
			TeamID:           "",
			OneMessageByTeam: false,
			rawQuery:         "cat rat",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000001",
				},
			},
		},
		{
			name:             "Success/RatAndCatUser3",
			UserID:           "00000000-0000-0000-0000-000000000003",
			TeamID:           "",
			OneMessageByTeam: false,
			rawQuery:         "cat rat",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000002",
					MessageID: "00000000-0000-0000-0000-000000000004",
				},
			},
		},
		{
			name:             "Success/RatAndCatTeam1",
			UserID:           "",
			TeamID:           "00000000-0000-0000-0000-000000000001",
			OneMessageByTeam: false,
			rawQuery:         "cat rat",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000001",
				},
			},
		},
		{
			name:             "Success/RatUser1",
			UserID:           "00000000-0000-0000-0000-000000000001",
			TeamID:           "",
			OneMessageByTeam: false,
			rawQuery:         "rat",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000002",
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
			name:             "Success/RatUser1OneMessage",
			UserID:           "00000000-0000-0000-0000-000000000001",
			TeamID:           "",
			OneMessageByTeam: true,
			rawQuery:         "rat",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000002",
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
			name:             "Success/EmptyUser",
			UserID:           "",
			TeamID:           "00000000-0000-0000-0000-000000000001",
			OneMessageByTeam: false,
			rawQuery:         "",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000007",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000006",
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
		{
			name:             "Success/EmptyTeam",
			UserID:           "00000000-0000-0000-0000-000000000001",
			TeamID:           "",
			OneMessageByTeam: false,
			rawQuery:         "",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000007",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000006",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000002",
					MessageID: "00000000-0000-0000-0000-000000000005",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000002",
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
		{
			name:             "Success/EmptyTeamOneMessage",
			UserID:           "00000000-0000-0000-0000-000000000001",
			TeamID:           "",
			OneMessageByTeam: true,
			rawQuery:         "",
			expect: []*entities.Message{
				{
					TeamID:    "00000000-0000-0000-0000-000000000001",
					MessageID: "00000000-0000-0000-0000-000000000007",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000002",
					MessageID: "00000000-0000-0000-0000-000000000005",
				},
				{
					TeamID:    "00000000-0000-0000-0000-000000000002",
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
			},
		},
	}

	stx := BeginTX(db, searchMessageFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {

			repo := dao.NewSearchMessagesRepository(stx)
			Message, err := repo.SearchMessages(context.TODO(), &dao.SearchMessageData{
				UserID:           tt.UserID,
				TeamID:           tt.TeamID,
				RawQuery:         tt.rawQuery,
				Limit:            1000,
				Offset:           0,
				OneMessageByTeam: tt.OneMessageByTeam,
			})

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
