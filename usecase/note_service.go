package usecase

import (
	"errors"
	"time"

	"my-note/domain"
	"my-note/repository"

	"github.com/google/uuid"
)

type NoteService struct {
	repo repository.NoteRepository
}

func NewNoteService(r repository.NoteRepository) *NoteService {
	return &NoteService{repo: r}
}

func (s *NoteService) CreateNote(title, content string) (domain.Note, error) {
	n := domain.Note{
		ID:        uuid.NewString(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	return s.repo.Create(n)
}

func (s *NoteService) GetNote(id string) (domain.Note, error) {
	n, ok, err := s.repo.Get(id)
	if err != nil {
		return domain.Note{}, err
	}
	if !ok {
		return domain.Note{}, errors.New("note not found")
	}
	return n, nil
}

func (s *NoteService) ListNotes() ([]domain.Note, error) {
	return s.repo.List()
}
