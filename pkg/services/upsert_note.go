package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/models"
)

type UpsertNoteService interface {
	Exec(ctx context.Context, note *models.UpsertNote) (*models.Note, error)
}
type upsertNoteServiceImpl struct {
	updateNoteRepository dao.UpdateNoteRepository
	createNoteRepository dao.CreateNoteRepository
	deleteNoteRepository dao.DeleteNoteRepository
}

func (s *upsertNoteServiceImpl) Exec(ctx context.Context, note *models.UpsertNote) (*models.Note, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(note); err != nil {
		return nil, errors.Join(ErrInvalidNoteUpdate, err)
	}

	targetName := note.TargetName + " " + note.PublicIdentifier

	// Delete note if content is empty
	if note.Content == "" {
		return nil, s.deleteNoteRepository.DeleteNote(ctx, note.AuthorID, note.NoteID)
	}

	// Attempt to create a note.
	createdNote, err := s.createNoteRepository.CreateNote(
		ctx,
		note.AuthorID,
		note.NoteID,
		&dao.CreateNoteData{
			NoteContent: note.Content,
			TargetName:  targetName,
		})

	// Note was successfully created.
	if err == nil {
		return &models.Note{
			AuthorID:   createdNote.AuthorID,
			NoteID:     createdNote.NoteID,
			Content:    createdNote.Content,
			TargetName: createdNote.TargetName,
		}, nil
	}

	if !errors.Is(err, dao.ErrNoteAlreadyExists) {
		return nil, err
	}

	updatedNote, err := s.updateNoteRepository.UpdateNote(
		ctx,
		note.AuthorID,
		note.NoteID,
		&dao.UpdateNoteData{
			NoteContent: note.Content,
			TargetName:  targetName,
		})

	if err != nil {
		return nil, err
	}

	return &models.Note{
		AuthorID:   updatedNote.AuthorID,
		NoteID:     updatedNote.NoteID,
		Content:    updatedNote.Content,
		TargetName: updatedNote.TargetName,
	}, nil
}

func NewUpsertNoteService(
	updateNoteRepository dao.UpdateNoteRepository,
	createNoteRepository dao.CreateNoteRepository,
	deleteNoteRepository dao.DeleteNoteRepository,
) UpsertNoteService {
	return &upsertNoteServiceImpl{
		updateNoteRepository: updateNoteRepository,
		createNoteRepository: createNoteRepository,
		deleteNoteRepository: deleteNoteRepository,
	}
}
