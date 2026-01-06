package persistence

import (
	"context"
	"fmt"
	"sync"

	"github.com/etsrc/goprod/internal/domain"
)

type InMemoryBookmarkRepository struct {
	mu        sync.RWMutex
	bookmarks map[string]*domain.Bookmark
}

func NewInMemoryBookmarkRepository() *InMemoryBookmarkRepository {
	return &InMemoryBookmarkRepository{
		bookmarks: make(map[string]*domain.Bookmark),
	}
}

func (r *InMemoryBookmarkRepository) Create(_ context.Context, b *domain.Bookmark) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.bookmarks[b.ID]; exists {
		// Given the architecture, the service layer generates UUIDs, so this shouldn't happen.
		return fmt.Errorf("persistence.InMemoryBookmarkRepository.Create: bookmark with ID %s already exists", b.ID)
	}
	r.bookmarks[b.ID] = b
	return nil
}

func (r *InMemoryBookmarkRepository) GetByID(_ context.Context, id string) (*domain.Bookmark, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bookmark, ok := r.bookmarks[id]
	if !ok {
		return nil, domain.ErrBookmarkNotFound
	}
	return bookmark, nil
}

func (r *InMemoryBookmarkRepository) GetAll(_ context.Context) ([]*domain.Bookmark, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allBookmarks := make([]*domain.Bookmark, 0, len(r.bookmarks))
	for _, bookmark := range r.bookmarks {
		allBookmarks = append(allBookmarks, bookmark)
	}
	return allBookmarks, nil
}

func (r *InMemoryBookmarkRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.bookmarks[id]; !ok {
		return domain.ErrBookmarkNotFound
	}
	delete(r.bookmarks, id)
	return nil
}
