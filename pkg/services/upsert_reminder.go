package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/models"
)

type UpsertReminderService interface {
	Exec(ctx context.Context, reminder *models.UpsertReminder) (*models.Reminder, error)
}
type upsertReminderServiceImpl struct {
	updateReminderRepository dao.UpdateReminderRepository
	createReminderRepository dao.CreateReminderRepository
	deleteReminderRepository dao.DeleteReminderRepository
}

func (s *upsertReminderServiceImpl) Exec(ctx context.Context, reminder *models.UpsertReminder) (*models.Reminder, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(reminder); err != nil {
		return nil, errors.Join(ErrInvalidReminderUpdate, err)
	}

	targetName := reminder.TargetName + " " + reminder.PublicIdentifier

	// Delete message if content is empty
	if reminder.Content == "" {
		return nil, s.deleteReminderRepository.DeleteReminder(ctx, reminder.AuthorID, reminder.ReminderID)
	}

	// Attempt to create a message.
	createdReminder, err := s.createReminderRepository.CreateReminder(
		ctx,
		reminder.AuthorID,
		reminder.ReminderID,
		&dao.CreateReminderData{
			ReminderContent: reminder.Content,
			TargetName:      targetName,
			ExpiredAt:       reminder.ExpiredAt,
		})

	// Reminder was successfully created.
	if err == nil {
		return &models.Reminder{
			AuthorID:   createdReminder.AuthorID,
			ReminderID: createdReminder.ReminderID,
			Content:    createdReminder.Content,
			TargetName: createdReminder.TargetName,
			ExpiredAt:  createdReminder.ExpiredAt,
		}, nil
	}

	if !errors.Is(err, dao.ErrReminderAlreadyExists) {
		return nil, err
	}

	updatedReminder, err := s.updateReminderRepository.UpdateReminder(
		ctx,
		reminder.AuthorID,
		reminder.ReminderID,
		&dao.UpdateReminderData{
			ReminderContent: reminder.Content,
			TargetName:      targetName,
			ExpiredAt:       reminder.ExpiredAt,
		})

	if err != nil {
		return nil, err
	}

	return &models.Reminder{
		AuthorID:   updatedReminder.AuthorID,
		ReminderID: updatedReminder.ReminderID,
		Content:    updatedReminder.Content,
		TargetName: updatedReminder.TargetName,
		ExpiredAt:  updatedReminder.ExpiredAt,
	}, nil
}

func NewUpsertReminderService(
	updateReminderRepository dao.UpdateReminderRepository,
	createReminderRepository dao.CreateReminderRepository,
	deleteReminderRepository dao.DeleteReminderRepository,
) UpsertReminderService {
	return &upsertReminderServiceImpl{
		updateReminderRepository: updateReminderRepository,
		createReminderRepository: createReminderRepository,
		deleteReminderRepository: deleteReminderRepository,
	}
}
