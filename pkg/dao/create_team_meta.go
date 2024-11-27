package dao

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type CreateTeamMetaRepository interface {
	CreateTeamMeta(ctx context.Context, teamID string, userID string) (*entities.TeamMeta, error)
}

type createTeamMetaRepositoryImpl struct {
	db bun.IDB
}

func (r *createTeamMetaRepositoryImpl) CreateTeamMeta(ctx context.Context, teamID string, userID string) (*entities.TeamMeta, error) {
	teamMeta := &entities.TeamMeta{
		TeamID: lo.ToPtr(uuid.MustParse(teamID)),
		UserID: lo.ToPtr(uuid.MustParse(userID)),
	}

	if _, err := r.db.NewInsert().Model(teamMeta).Returning("*").Exec(ctx); err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			return nil, ErrNoteAlreadyExists
		}

		return nil, err
	}

	return teamMeta, nil
}

func NewCreateTeamMetaRepository(db bun.IDB) CreateTeamMetaRepository {
	return &createTeamMetaRepositoryImpl{
		db: db,
	}
}
