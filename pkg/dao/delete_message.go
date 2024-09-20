package dao

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type DeleteMessageRepository interface {
	DeleteMessage(ctx context.Context, teamID string, messageID string) error
}

type deleteMessageRepositoryImpl struct {
	db bun.IDB
}

func (r *deleteMessageRepositoryImpl) DeleteMessage(ctx context.Context, teamID string, messageID string) error {
	_, err := r.db.NewDelete().
		Model(&entities.Message{}).
		Where("team_id = ?", teamID).
		Where("message_id = ?", messageID).
		Exec(ctx)

	return err
}

func NewDeleteMessageRepository(db bun.IDB) DeleteMessageRepository {
	return &deleteMessageRepositoryImpl{
		db: db,
	}
}
