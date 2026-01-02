package domain

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"
)

type BookmarkRepository interface {
	Create(ctx context.Context, b *Bookmark) error
	GetByID(ctx context.Context, id string) (*Bookmark, error)
	GetAll(ctx context.Context) ([]*Bookmark, error)
	Delete(ctx context.Context, id string) error
}

var (
	ErrBookmarkNotFound = errors.New("bookmark not found")
	ErrInvalidURL       = errors.New("the provided URL is invalid")
	ErrTitleTooShort    = errors.New("title must be at least 3 characters")
)

type Bookmark struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *Bookmark) Validate() error {
	if b.Title == "" || len(strings.TrimSpace(b.Title)) < 3 {
		return ErrTitleTooShort
	}

	_, err := url.ParseRequestURI(b.URL)
	if err != nil {
		return ErrInvalidURL
	}

	return nil
}
