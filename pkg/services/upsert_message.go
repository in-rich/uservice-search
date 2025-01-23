package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/samber/lo"
	"time"
)

type UpsertMessageService interface {
	Exec(ctx context.Context, message *models.UpsertMessage) (*models.Message, error)
}
type upsertMessageServiceImpl struct {
	updateMessageRepository dao.UpdateMessageRepository
	createMessageRepository dao.CreateMessageRepository
	deleteMessageRepository dao.DeleteMessageRepository
}

func (s *upsertMessageServiceImpl) Exec(ctx context.Context, message *models.UpsertMessage) (*models.Message, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(message); err != nil {
		return nil, errors.Join(ErrInvalidMessageUpdate, err)
	}

	targetName := message.TargetName + " " + message.PublicIdentifier

	// Delete message if content is empty
	if message.Content == "" {
		return nil, s.deleteMessageRepository.DeleteMessage(ctx, message.TeamID, message.MessageID)
	}

	if message.MessageID == "" {
		return nil, ErrInvalidMessageUpdate
	}

	updatedAt, _ := lo.Coalesce(lo.FromPtr(message.UpdatedAt), time.Now())

	// Attempt to create a message.
	createdMessage, err := s.createMessageRepository.CreateMessage(
		ctx,
		message.TeamID,
		message.MessageID,
		&dao.CreateMessageData{
			MessageContent: message.Content,
			TargetName:     targetName,
			UpdatedAt:      updatedAt,
		})

	// Note was successfully created.
	if err == nil {
		return &models.Message{
			TeamID:     createdMessage.TeamID,
			MessageID:  createdMessage.MessageID,
			Content:    createdMessage.Content,
			TargetName: createdMessage.TargetName,
		}, nil
	}

	if !errors.Is(err, dao.ErrMessageAlreadyExists) {
		return nil, err
	}

	updatedMessage, err := s.updateMessageRepository.UpdateMessage(
		ctx,
		message.TeamID,
		message.MessageID,
		&dao.UpdateMessageData{
			MessageContent: message.Content,
			TargetName:     targetName,
			UpdatedAt:      updatedAt,
		})

	if err != nil {
		return nil, err
	}

	return &models.Message{
		TeamID:     updatedMessage.TeamID,
		MessageID:  updatedMessage.MessageID,
		Content:    updatedMessage.Content,
		TargetName: updatedMessage.TargetName,
	}, nil
}

func NewUpsertMessageService(
	updateMessageRepository dao.UpdateMessageRepository,
	createMessageRepository dao.CreateMessageRepository,
	deleteMessageRepository dao.DeleteMessageRepository,
) UpsertMessageService {
	return &upsertMessageServiceImpl{
		updateMessageRepository: updateMessageRepository,
		createMessageRepository: createMessageRepository,
		deleteMessageRepository: deleteMessageRepository,
	}
}
