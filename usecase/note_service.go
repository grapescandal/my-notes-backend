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

func (s *NoteService) DeleteNote(id string) error {
	ok, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("note not found")
	}
	return nil
}

func (s *NoteService) SaveNote(n domain.Note) (domain.Note, bool, error) {
	// returns (note, created, error)
	if n.ID == "" {
		n.ID = uuid.NewString()
		n.CreatedAt = time.Now().UTC()
		nn, err := s.repo.Create(n)
		return nn, true, err
	}

	// Check if exists
	existing, ok, err := s.repo.Get(n.ID)
	if err != nil {
		return domain.Note{}, false, err
	}
	if ok {
		// preserve CreatedAt from existing
		n.CreatedAt = existing.CreatedAt
		nn, err := s.repo.Update(n)
		return nn, false, err
	}

	// Not found: create new with provided ID
	n.CreatedAt = time.Now().UTC()
	nn, err := s.repo.Create(n)
	return nn, true, err
}
