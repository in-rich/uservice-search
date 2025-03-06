package dao

import (
	"context"
	"github.com/google/uuid"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
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
		UserID: userID,
	}

	_, err := r.db.NewInsert().
		Model(teamMeta).
		On("CONFLICT (team_id, user_id) DO NOTHING").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return teamMeta, nil
}

func NewCreateTeamMetaRepository(db bun.IDB) CreateTeamMetaRepository {
	return &createTeamMetaRepositoryImpl{
		db: db,
	}
}
