package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/models"
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

	searchNotes, err := s.searchNotesRepository.SearchNotes(ctx, note.Limit, note.Offset, note.RawQuery, note.AuthorID)
	if err != nil {
		if errors.Is(err, dao.ErrNoteNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}

	modelsNotes := make([]*models.Note, 0)

	for _, tt := range searchNotes {
		modelsNotes = append(modelsNotes, &models.Note{
			AuthorID: tt.AuthorID,
			NoteID:   tt.NoteID,
		})
	}

	return modelsNotes, nil
}

func NewSearchNotesService(searchNotesRepository dao.SearchNotesRepository) SearchNotesService {
	return &searchNotesServiceImpl{
		searchNotesRepository: searchNotesRepository,
	}
}
