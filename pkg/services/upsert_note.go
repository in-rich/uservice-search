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

	// Delete message if content is empty
	if note.Content == "" {
		return nil, s.deleteNoteRepository.DeleteNote(ctx, note.AuthorID, note.NoteID)
	}

	updatedAt, _ := lo.Coalesce(lo.FromPtr(note.UpdatedAt), time.Now())

	// Attempt to create a message.
	createdNote, err := s.createNoteRepository.CreateNote(
		ctx,
		note.AuthorID,
		note.NoteID,
		&dao.CreateNoteData{
			NoteContent: note.Content,
			TargetName:  targetName,
			UpdatedAt:   updatedAt,
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
			UpdatedAt:   updatedAt,
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
