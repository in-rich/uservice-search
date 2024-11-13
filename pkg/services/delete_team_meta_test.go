package services_test

import (
	"context"
	daomocks "github.com/in-rich/uservice-search/pkg/dao/mocks"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteTeam(t *testing.T) {
	testData := []struct {
		name string

		teamID string
		userID string

		shouldCallDeleteTeamMetaRepository bool
		deleteTeamMetaRepositoryErr        error

		shouldCallDeleteTeamMetaMemberRepository bool
		deleteTeamMetaMemberRepositoryErr        error

		expectErr error
	}{
		{
			name:                                     "DeleteTeamMetaMember",
			teamID:                                   "00000000-0000-0000-0000-000000000001",
			userID:                                   "00000000-0000-0000-0000-000000000001",
			shouldCallDeleteTeamMetaRepository:       false,
			shouldCallDeleteTeamMetaMemberRepository: true,
		},
		{
			name:                                     "DeleteTeamMeta",
			teamID:                                   "00000000-0000-0000-0000-000000000001",
			userID:                                   "",
			shouldCallDeleteTeamMetaRepository:       true,
			shouldCallDeleteTeamMetaMemberRepository: false,
		},
		{
			name:                                     "DeleteTeamMetaError",
			teamID:                                   "00000000-0000-0000-0000-000000000001",
			userID:                                   "00000000-0000-0000-0000-000000000001",
			shouldCallDeleteTeamMetaRepository:       false,
			shouldCallDeleteTeamMetaMemberRepository: true,
			deleteTeamMetaMemberRepositoryErr:        FooErr,
			expectErr:                                FooErr,
		},
		{
			name:      "InvalidData",
			teamID:    "",
			userID:    "00000000-0000-0000-0000-000000000001",
			expectErr: services.ErrInvalidTeamMetaDelete,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			deleteTeamMetaRepository := daomocks.NewMockDeleteTeamMetaRepository(t)
			deleteTeamMetaMemberRepository := daomocks.NewMockDeleteTeamMetaMemberRepository(t)

			if tt.shouldCallDeleteTeamMetaRepository {
				deleteTeamMetaRepository.On("DeleteTeamMeta", context.TODO(), tt.teamID).Return(tt.deleteTeamMetaRepositoryErr)
			}

			if tt.shouldCallDeleteTeamMetaMemberRepository {
				deleteTeamMetaMemberRepository.On("DeleteTeamMetaMember", context.TODO(), tt.teamID, tt.userID).Return(tt.deleteTeamMetaMemberRepositoryErr)
			}

			service := services.NewDeleteTeamMetaService(deleteTeamMetaRepository, deleteTeamMetaMemberRepository)

			err := service.Exec(context.TODO(), &models.DeleteTeamMeta{TeamID: tt.teamID, UserID: tt.userID})

			require.ErrorIs(t, err, tt.expectErr)

			deleteTeamMetaRepository.AssertExpectations(t)
		})
	}
}
