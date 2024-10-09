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

type SearchMessagesService interface {
	Exec(ctx context.Context, note *models.SearchMessages) ([]*models.Message, error)
}

type searchMessagesServiceImpl struct {
	searchMessageRepository dao.SearchMessagesRepository
}

func (s *searchMessagesServiceImpl) Exec(ctx context.Context, message *models.SearchMessages) ([]*models.Message, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(message); err != nil {
		return nil, errors.Join(ErrInvalidMessageSearch, err)
	}

	searchMessages, err := s.searchMessageRepository.SearchMessages(ctx, message.TeamID, message.RawQuery, message.Limit, message.Offset)
	if err != nil {
		return nil, err
	}

	modelsMessages := lo.Map(searchMessages, func(message *entities.Message, _ int) *models.Message {
		return &models.Message{
			TeamID:    message.TeamID,
			MessageID: message.MessageID,
		}
	})

	return modelsMessages, nil
}

func NewSearchMessagesService(searchMessageRepository dao.SearchMessagesRepository) SearchMessagesService {
	return &searchMessagesServiceImpl{
		searchMessageRepository: searchMessageRepository,
	}
}
