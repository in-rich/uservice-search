package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/entities"
	"github.com/in-rich/uservice-search/pkg/models"
	"github.com/samber/lo"
)

type SearchRemindersService interface {
	Exec(ctx context.Context, reminder *models.SearchReminders) ([]*models.Reminder, error)
}

type searchRemindersServiceImpl struct {
	searchRemindersRepository dao.SearchRemindersRepository
}

func (s *searchRemindersServiceImpl) Exec(ctx context.Context, reminder *models.SearchReminders) ([]*models.Reminder, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(reminder); err != nil {
		return nil, errors.Join(ErrInvalidReminderSearch, err)
	}

	searchReminders, err := s.searchRemindersRepository.SearchReminders(ctx, reminder.AuthorID, reminder.RawQuery, reminder.Limit, reminder.Offset)
	if err != nil {
		return nil, err
	}

	modelsReminders := lo.Map(searchReminders, func(reminder *entities.Reminder, _ int) *models.Reminder {
		return &models.Reminder{
			AuthorID:   reminder.AuthorID,
			ReminderID: reminder.ReminderID,
		}
	})
	return modelsReminders, nil
}

func NewSearchRemindersService(searchRemindersRepository dao.SearchRemindersRepository) SearchRemindersService {
	return &searchRemindersServiceImpl{
		searchRemindersRepository: searchRemindersRepository,
	}
}
