package dao

import (
	"context"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type DeleteReminderRepository interface {
	DeleteReminder(ctx context.Context, authorID string, reminderID string) error
}

type deleteReminderRepositoryImpl struct {
	db bun.IDB
}

func (r *deleteReminderRepositoryImpl) DeleteReminder(ctx context.Context, authorID string, reminderID string) error {
	_, err := r.db.NewDelete().
		Model(&entities.Reminder{}).
		Where("author_id = ?", authorID).
		Where("reminder_id = ?", reminderID).
		Exec(ctx)

	return err
}

func NewDeleteReminderRepository(db bun.IDB) DeleteReminderRepository {
	return &deleteReminderRepositoryImpl{
		db: db,
	}
}
