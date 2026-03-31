package memory

import (
	"sync"

	"my-note/domain"
)

type MemoryRepo struct {
	mu    sync.RWMutex
	notes map[string]domain.Note
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{notes: make(map[string]domain.Note)}
}

func (r *MemoryRepo) Create(n domain.Note) (domain.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notes[n.ID] = n
	return n, nil
}

func (r *MemoryRepo) Get(id string) (domain.Note, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.notes[id]
	return n, ok, nil
}

func (r *MemoryRepo) List() ([]domain.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Note, 0, len(r.notes))
	for _, v := range r.notes {
		out = append(out, v)
	}
	return out, nil
}
