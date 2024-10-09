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

type SearchNotesService interface {
	Exec(ctx context.Context, note *models.SearchNotes) ([]*models.Note, error)
}

type searchNotesServiceImpl struct {
	searchNotesRepository dao.SearchNotesRepository
}

func (s *searchNotesServiceImpl) Exec(ctx context.Context, note *models.SearchNotes) ([]*models.Note, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(note); err != nil {
		return nil, errors.Join(ErrInvalidNoteSearch, err)
	}

	searchNotes, err := s.searchNotesRepository.SearchNotes(ctx, note.AuthorID, note.RawQuery, note.Limit, note.Offset)
	if err != nil {
		return nil, err
	}

	modelsNotes := lo.Map(searchNotes, func(note *entities.Note, _ int) *models.Note {
		return &models.Note{
			AuthorID: note.AuthorID,
			NoteID:   note.NoteID,
		}
	})

	return modelsNotes, nil
}

func NewSearchNotesService(searchNotesRepository dao.SearchNotesRepository) SearchNotesService {
	return &searchNotesServiceImpl{
		searchNotesRepository: searchNotesRepository,
	}
}
