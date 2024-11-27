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

var deleteTeamMetaFixtures = []*entities.TeamMeta{
	{
		TeamID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		UserID: "00000000-0000-0000-0000-000000000001",
	},
}

func TestDeleteTeamMeta(t *testing.T) {
	db := OpenDB()
	defer CloseDB(db)

	testData := []struct {
		name string

		teamID string
		userID string

		expectErr error
	}{
		{
			name:   "DeleteTeamMeta",
			teamID: "00000000-0000-0000-0000-000000000001",
		},
		{
			// Still a success because this method if forgiving.
			name:   "TeamMetaNotFound",
			teamID: "00000000-0000-0000-0000-000000000001",
		},
	}

	stx := BeginTX(db, deleteTeamMetaFixtures)
	defer RollbackTX(stx)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tx := BeginTX[interface{}](stx, nil)
			defer RollbackTX(tx)

			repo := dao.NewDeleteTeamMetaRepository(tx)
			err := repo.DeleteTeamMeta(context.TODO(), tt.teamID)

			require.ErrorIs(t, err, tt.expectErr)
		})
	}
}
