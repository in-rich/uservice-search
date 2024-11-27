package dao

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type DeleteTeamMetaMemberRepository interface {
	DeleteTeamMetaMember(ctx context.Context, teamID string, userID string) error
}

type deleteTeamMetaMemberRepositoryImpl struct {
	db bun.IDB
}

func (r *deleteTeamMetaMemberRepositoryImpl) DeleteTeamMetaMember(ctx context.Context, teamID string, userID string) error {
	_, err := r.db.NewDelete().
		Model(&entities.TeamMeta{}).
		Where("team_id = ?", teamID).
		Where("user_id = ?", userID).
		Exec(ctx)

	return err
}

func NewDeleteTeamMetaMemberRepository(db bun.IDB) DeleteTeamMetaMemberRepository {
	return &deleteTeamMetaMemberRepositoryImpl{
		db: db,
	}
}
