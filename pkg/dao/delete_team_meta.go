package dao

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type DeleteTeamMetaRepository interface {
	DeleteTeamMeta(ctx context.Context, teamID string) error
}

type deleteTeamMetaRepositoryImpl struct {
	db bun.IDB
}

func (r *deleteTeamMetaRepositoryImpl) DeleteTeamMeta(ctx context.Context, teamID string) error {
	_, err := r.db.NewDelete().
		Model(&entities.TeamMeta{}).
		Where("team_id = ?", teamID).
		Exec(ctx)

	return err
}

func NewDeleteTeamMetaRepository(db bun.IDB) DeleteTeamMetaRepository {
	return &deleteTeamMetaRepositoryImpl{
		db: db,
	}
}
