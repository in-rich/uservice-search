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

var deleteTeamMetaMemberFixtures = []*entities.TeamMeta{
	{
		TeamID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		UserID: "00000000-0000-0000-0000-000000000001",
	},
}

func TestDeleteTeamMetaMember(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		teamID string
		userID string

		expectErr error
	}{
		{
			name:   "DeleteTeamMetaMember",
			teamID: "00000000-0000-0000-0000-000000000001",
			userID: "00000000-0000-0000-0000-000000000001",
		},
		{
			// Still a success because this method if forgiving.
			name:   "TeamMetaNotFound",
			teamID: "00000000-0000-0000-0000-000000000001",
			userID: "00000000-0000-0000-0000-000000000002",
		},
	}

	stx := BeginTX(db, deleteTeamMetaMemberFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewDeleteTeamMetaMemberRepository(tx)
			err := repo.DeleteTeamMetaMember(context.TODO(), tt.teamID, tt.userID)

			require.ErrorIs(t, err, tt.expectErr)
		})
	}
}
