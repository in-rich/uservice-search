package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/models"
)

type CreateTeamMetaService interface {
	Exec(ctx context.Context, teamMeta *models.TeamMeta) (*models.TeamMeta, error)
}
type createTeamMetaServiceImpl struct {
	createTeamMetaRepository dao.CreateTeamMetaRepository
}

func (s *createTeamMetaServiceImpl) Exec(ctx context.Context, teamMeta *models.TeamMeta) (*models.TeamMeta, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(teamMeta); err != nil {
		return nil, errors.Join(ErrInvalidTeamMetaCreate, err)
	}

	// Attempt to create a message.
	createdTeamMeta, err := s.createTeamMetaRepository.CreateTeamMeta(
		ctx,
		teamMeta.TeamID,
		teamMeta.UserID,
	)

	if err != nil {
		return nil, err
	}

	modelsTeamMeta := &models.TeamMeta{
		TeamID: createdTeamMeta.TeamID.String(),
		UserID: createdTeamMeta.UserID,
	}

	return modelsTeamMeta, nil
}

func NewCreateTeamMetaService(
	createTeamMetaRepository dao.CreateTeamMetaRepository,
) CreateTeamMetaService {
	return &createTeamMetaServiceImpl{
		createTeamMetaRepository: createTeamMetaRepository,
	}
}
