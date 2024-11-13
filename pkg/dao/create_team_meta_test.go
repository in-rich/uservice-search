package dao_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
)

var createTeamMetaFixtures = []*entities.TeamMeta{
	{
		TeamID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		UserID: "00000000-0000-0000-0000-000000000001",
	},
}

func TestCreateTeamMeta(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		teamID string
		UserID string

		expect    *entities.TeamMeta
		expectErr error
	}{
		{
			name:   "Success",
			teamID: "00000000-0000-0000-0000-000000000001",
			UserID: "00000000-0000-0000-0000-000000000002",
			expect: &entities.TeamMeta{
				TeamID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				UserID: "00000000-0000-0000-0000-000000000002",
			},
		},
		{
			name:      "Error/TeamMetaAlreadyExists",
			teamID:    "00000000-0000-0000-0000-000000000001",
			UserID:    "00000000-0000-0000-0000-000000000001",
			expectErr: dao.ErrTeamMetaAlreadyExists,
		},
		{
			name:   "Success/SameUserDifferentTeam",
			teamID: "00000000-0000-0000-0000-000000000002",
			UserID: "00000000-0000-0000-0000-000000000001",
			expect: &entities.TeamMeta{
				TeamID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
				UserID: "00000000-0000-0000-0000-000000000001",
			},
		},
	}

	stx := BeginTX(db, createTeamMetaFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewCreateTeamMetaRepository(tx)
			message, err := repo.CreateTeamMeta(context.TODO(), tt.teamID, tt.UserID)

			if message != nil {
				// Since ID is random, nullify if for comparison.
				message.ID = nil
			}

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, message)
		})
	}
}
