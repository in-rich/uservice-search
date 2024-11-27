package dao

import (
	"context"
	"errors"

	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type CreateMessageData struct {
	MessageContent string
	TargetName     string
}

type CreateMessageRepository interface {
	CreateMessage(ctx context.Context, teamID string, messageID string, data *CreateMessageData) (*entities.Message, error)
}

type createMessageRepositoryImpl struct {
	db bun.IDB
}

func (r *createMessageRepositoryImpl) CreateMessage(ctx context.Context, teamID string, messageID string, data *CreateMessageData) (*entities.Message, error) {
	message := &entities.Message{
		TeamID:     teamID,
		MessageID:  messageID,
		Content:    data.MessageContent,
		TargetName: data.TargetName,
	}

	_, err := message.BeforeCreate(r.db.NewInsert().Model(message)).Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			return nil, ErrMessageAlreadyExists
		}

		return nil, err
	}

	return message, nil
}

func NewCreateMessageRepository(db bun.IDB) CreateMessageRepository {
	return &createMessageRepositoryImpl{
		db: db,
	}
}
