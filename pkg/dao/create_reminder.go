package dao

import (
	"context"
	"errors"
	"time"

	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type CreateReminderData struct {
	ReminderContent string
	TargetName      string
	ExpiredAt       *time.Time
}

type CreateReminderRepository interface {
	CreateReminder(ctx context.Context, authorID string, reminderID string, data *CreateReminderData) (*entities.Reminder, error)
}

type createReminderRepositoryImpl struct {
	db bun.IDB
}

func (r *createReminderRepositoryImpl) CreateReminder(ctx context.Context, authorID string, reminderID string, data *CreateReminderData) (*entities.Reminder, error) {
	reminder := &entities.Reminder{
		AuthorID:   authorID,
		ReminderID: reminderID,
		Content:    data.ReminderContent,
		TargetName: data.TargetName,
		ExpiredAt:  data.ExpiredAt,
	}

	_, err := reminder.BeforeCreate(r.db.NewInsert().Model(reminder)).Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			return nil, ErrReminderAlreadyExists
		}

		return nil, err
	}

	return reminder, nil
}

func NewCreateReminderRepository(db bun.IDB) CreateReminderRepository {
	return &createReminderRepositoryImpl{
		db: db,
	}
}
