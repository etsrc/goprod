package service

import (
	"context"
	"fmt"
	"time"

	"github.com/etsrc/goprod/internal/domain"
	"github.com/google/uuid"
)

type BookmarkService interface {
	Create(ctx context.Context, b *domain.Bookmark) error
	GetByID(ctx context.Context, id string) (*domain.Bookmark, error)
	List(ctx context.Context) ([]*domain.Bookmark, error)
	Delete(ctx context.Context, id string) error
}

type bookmarkService struct {
	repo domain.BookmarkRepository
}

func NewBookmarkService(repo domain.BookmarkRepository) BookmarkService {
	return &bookmarkService{
		repo: repo,
	}
}

func (s *bookmarkService) Create(ctx context.Context, b *domain.Bookmark) error {
	b.ID = uuid.NewString()
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	if err := b.Validate(); err != nil {
		return fmt.Errorf("service.Create: %w", err)
	}

	if err := s.repo.Create(ctx, b); err != nil {
		return fmt.Errorf("service.Create: failed to save: %w", err)
	}

	return nil
}

func (s *bookmarkService) GetByID(ctx context.Context, id string) (*domain.Bookmark, error) {
	if id == "" {
		return nil, fmt.Errorf("service.GetByID: id is required")
	}

	bookmark, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID: %w", err)
	}

	return bookmark, nil
}

func (s *bookmarkService) List(ctx context.Context) ([]*domain.Bookmark, error) {
	bookmarks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.List: %w", err)
	}

	return bookmarks, nil
}

func (s *bookmarkService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("service.Delete: id is required")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service.Delete: %w", err)
	}

	return nil
}
