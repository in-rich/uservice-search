package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/models"
)

type DeleteTeamMetaService interface {
	Exec(ctx context.Context, teamMeta *models.DeleteTeamMeta) error
}
type deleteTeamMetaServiceImpl struct {
	deleteTeamMetaRepository       dao.DeleteTeamMetaRepository
	deleteTeamMetaMemberRepository dao.DeleteTeamMetaMemberRepository
}

func (s *deleteTeamMetaServiceImpl) Exec(ctx context.Context, teamMeta *models.DeleteTeamMeta) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(teamMeta); err != nil {
		return errors.Join(ErrInvalidTeamMetaDelete, err)
	}

	// If the userID is empty we delete all the user associated with the teamID.
	if teamMeta.UserID == "" {
		return s.deleteTeamMetaRepository.DeleteTeamMeta(
			ctx,
			teamMeta.TeamID,
		)
	}
	// Else we only delete the userID provided for the teamID provided.
	return s.deleteTeamMetaMemberRepository.DeleteTeamMetaMember(
		ctx,
		teamMeta.TeamID,
		teamMeta.UserID,
	)
}

func NewDeleteTeamMetaService(
	deleteTeamMetaRepository dao.DeleteTeamMetaRepository,
	deleteTeamMetaMember dao.DeleteTeamMetaMemberRepository,
) DeleteTeamMetaService {
	return &deleteTeamMetaServiceImpl{
		deleteTeamMetaRepository:       deleteTeamMetaRepository,
		deleteTeamMetaMemberRepository: deleteTeamMetaMember,
	}
}
