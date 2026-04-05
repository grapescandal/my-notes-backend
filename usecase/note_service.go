package usecase

import (
	"errors"
	"sort"
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
	t := time.Now().UTC()
	id := uuid.NewString()
	n := domain.Note{
		ID:        id,
		Title:     title,
		Content:   content,
		CreatedAt: &t,
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
	notes, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	// Sort by UpdatedAt (most recent first), then CreatedAt (most recent first)
	getUpdated := func(n domain.Note) time.Time {
		if n.UpdatedAt != nil {
			return *n.UpdatedAt
		}
		return time.Time{}
	}

	getCreated := func(n domain.Note) time.Time {
		if n.CreatedAt != nil {
			return *n.CreatedAt
		}
		return time.Time{}
	}

	// stable sort using slice comparator
	sort.Slice(notes, func(i, j int) bool {
		ui := getUpdated(notes[i])
		uj := getUpdated(notes[j])
		if !ui.Equal(uj) {
			return ui.After(uj)
		}
		return getCreated(notes[i]).After(getCreated(notes[j]))
	})

	return notes, nil
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
		id := uuid.NewString()
		n.ID = id
		t := time.Now().UTC()
		n.CreatedAt = &t
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
		// stamp updated at
		t := time.Now().UTC()
		n.UpdatedAt = &t
		nn, err := s.repo.Update(n)
		return nn, false, err
	}

	// Not found: create new with provided ID
	t := time.Now().UTC()
	n.CreatedAt = &t
	nn, err := s.repo.Create(n)
	return nn, true, err
}
