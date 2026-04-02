package repository

import "my-note/domain"

type NoteRepository interface {
	Create(n domain.Note) (domain.Note, error)
	Get(id string) (domain.Note, bool, error)
	List() ([]domain.Note, error)
	Delete(id string) (bool, error)
	Update(n domain.Note) (domain.Note, error)
}
