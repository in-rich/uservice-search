package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/models"
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

	// Attempt to create a message.
	createdNote, err := s.createMessageRepository.CreateMessage(
		ctx,
		message.TeamID,
		message.MessageID,
		&dao.CreateMessageData{
			MessageContent: message.Content,
			TargetName:     targetName,
		})

	// Note was successfully created.
	if err == nil {
		return &models.Message{
			TeamID:     createdNote.TeamID,
			MessageID:  createdNote.MessageID,
			Content:    createdNote.Content,
			TargetName: createdNote.TargetName,
		}, nil
	}

	if !errors.Is(err, dao.ErrMessageAlreadyExists) {
		return nil, err
	}

	updatedNote, err := s.updateMessageRepository.UpdateMessage(
		ctx,
		message.TeamID,
		message.MessageID,
		&dao.UpdateMessageData{
			MessageContent: message.Content,
			TargetName:     targetName,
		})

	if err != nil {
		return nil, err
	}

	return &models.Message{
		TeamID:     updatedNote.TeamID,
		MessageID:  updatedNote.MessageID,
		Content:    updatedNote.Content,
		TargetName: updatedNote.TargetName,
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