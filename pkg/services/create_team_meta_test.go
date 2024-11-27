package services_test

import (
	"context"
	"github.com/google/uuid"
	daomocks "github.com/in-rich/uservice-search/pkg/dao/mocks"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/in-rich/uservice-search/pkg/services"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTeamMeta(t *testing.T) {
	testData := []struct {
		name string

		in *models.TeamMeta

		shouldCallCreateTeamMeta bool
		createTeamMetaResponse   *entities.TeamMeta
		createTeamMetaErr        error

		expect    *models.TeamMeta
		expectErr error
	}{
		{
			name: "CreateTeamMeta",
			in: &models.TeamMeta{
				TeamID: "00000000-0000-0000-0000-000000000001",
				UserID: "00000000-0000-0000-0000-000000000001",
			},
			shouldCallCreateTeamMeta: true,
			createTeamMetaResponse: &entities.TeamMeta{
				TeamID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				UserID: "00000000-0000-0000-0000-000000000001",
			},
			expect: &models.TeamMeta{
				TeamID: "00000000-0000-0000-0000-000000000001",
				UserID: "00000000-0000-0000-0000-000000000001",
			},
		},
		{
			name: "CreateTeamMetaDAOError",
			in: &models.TeamMeta{
				TeamID: "00000000-0000-0000-0000-000000000001",
				UserID: "00000000-0000-0000-0000-000000000001",
			},
			shouldCallCreateTeamMeta: true,
			createTeamMetaErr:        FooErr,
			expectErr:                FooErr,
		},
		{
			name: "InvalidData",
			in: &models.TeamMeta{
				TeamID: "",
				UserID: "00000000-0000-0000-0000-000000000001",
			},
			expectErr: services.ErrInvalidTeamMetaCreate,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			createTeamMetaRepository := daomocks.NewMockCreateTeamMetaRepository(t)

			if tt.shouldCallCreateTeamMeta {
				createTeamMetaRepository.
					On("CreateTeamMeta", context.TODO(), tt.in.TeamID, tt.in.UserID).
					Return(tt.createTeamMetaResponse, tt.createTeamMetaErr)
			}

			service := services.NewCreateTeamMetaService(createTeamMetaRepository)

			resp, err := service.Exec(context.TODO(), tt.in)

			require.ErrorIs(t, err, tt.expectErr)
			require.Equal(t, tt.expect, resp)

			createTeamMetaRepository.AssertExpectations(t)
		})
	}
}
