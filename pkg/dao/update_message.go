package dao

import (
	"context"

	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/uptrace/bun"
)

type UpdateMessageData struct {
	MessageContent string
	TargetName     string
}

type UpdateMessageRepository interface {
	UpdateMessage(ctx context.Context, teamID string, messageID string, data *UpdateMessageData) (*entities.Message, error)
}

type updateMessageRepositoryImpl struct {
	db bun.IDB
}

func (r *updateMessageRepositoryImpl) UpdateMessage(ctx context.Context, teamID string, messageID string, data *UpdateMessageData) (*entities.Message, error) {
	message := &entities.Message{
		Content:    data.MessageContent,
		TargetName: data.TargetName,
	}

	res, err := message.BeforeUpdate(r.db.NewUpdate().Model(message)).
		Where("team_id = ?", teamID).
		Where("message_id = ?", messageID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrMessageNotFound
	}

	return message, nil
}

func NewUpdateMessageRepository(db bun.IDB) UpdateMessageRepository {
	return &updateMessageRepositoryImpl{
		db: db,
	}
}
