package dao

import (
	"context"
	"github.com/samber/lo"
	"time"

	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type UpdateReminderData struct {
	ReminderContent string
	TargetName      string
	ExpiredAt       *time.Time
}

type UpdateReminderRepository interface {
	UpdateReminder(ctx context.Context, authorID string, reminderID string, data *UpdateReminderData) (*entities.Reminder, error)
}

type updateReminderRepositoryImpl struct {
	db bun.IDB
}

func (r *updateReminderRepositoryImpl) UpdateReminder(ctx context.Context, authorID string, reminderID string, data *UpdateReminderData) (*entities.Reminder, error) {
	reminder := &entities.Reminder{
		Content:    data.ReminderContent,
		TargetName: data.TargetName,
		UpdatedAt:  lo.ToPtr(time.Now()),
		ExpiredAt:  data.ExpiredAt,
	}

	res, err := reminder.BeforeUpdate(r.db.NewUpdate().Model(reminder)).
		Where("author_id = ?", authorID).
		Where("reminder_id = ?", reminderID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrReminderNotFound
	}

	return reminder, nil
}

func NewUpdateReminderRepository(db bun.IDB) UpdateReminderRepository {
	return &updateReminderRepositoryImpl{
		db: db,
	}
}
