package persistence

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/etsrc/goprod/internal/domain"
	"github.com/google/uuid"
)

// newTestRepo creates a new InMemoryBookmarkRepository for testing purposes.
func newTestRepo() *InMemoryBookmarkRepository {
	return NewInMemoryBookmarkRepository()
}

func TestInMemoryBookmarkRepository_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		bookmarks   []*domain.Bookmark // Bookmarks to pre-populate the repo
		newBookmark *domain.Bookmark   // Bookmark to create
		wantErr     error
	}{
		{
			name:      "Successfully Create Bookmark",
			bookmarks: []*domain.Bookmark{},
			newBookmark: &domain.Bookmark{
				ID:        uuid.New().String(),
				URL:       "https://example.com/new",
				Title:     "New Bookmark",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "Create Bookmark Already Exists",
			bookmarks: []*domain.Bookmark{
				{
					ID:        "existing-id",
					URL:       "https://example.com/existing",
					Title:     "Existing Bookmark",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			newBookmark: &domain.Bookmark{
				ID:        "existing-id",
				URL:       "https://example.com/duplicate",
				Title:     "Duplicate Bookmark",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: errors.New("persistence.InMemoryBookmarkRepository.Create: bookmark with ID existing-id already exists"),
		},
		{
			name: "Create another bookmark after one exists",
			bookmarks: []*domain.Bookmark{
				{
					ID:        "first-id",
					URL:       "https://example.com/first",
					Title:     "First Bookmark",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			newBookmark: &domain.Bookmark{
				ID:        uuid.New().String(),
				URL:       "https://example.com/second",
				Title:     "Second Bookmark",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable for parallel tests
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := newTestRepo()
			ctx := context.Background()

			// Pre-populate the repository
			for _, b := range tt.bookmarks {
				repo.Create(ctx, b) // Errors from pre-population are not tested here
			}

			err := repo.Create(ctx, tt.newBookmark)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Create() expected error, got nil")
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("Create() got error = %v, want error %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Create() unexpected error = %v", err)
				}
				// Verify the bookmark was actually created
				found, getErr := repo.GetByID(ctx, tt.newBookmark.ID)
				if getErr != nil {
					t.Errorf("Create() failed to retrieve created bookmark: %v", getErr)
				}
				if found.ID != tt.newBookmark.ID {
					t.Errorf("Create() created bookmark ID mismatch: got %s, want %s", found.ID, tt.newBookmark.ID)
				}
			}
		})
	}
}

func TestInMemoryBookmarkRepository_GetByID(t *testing.T) {
	t.Parallel()

	bookmark1 := &domain.Bookmark{
		ID:        "id-1",
		URL:       "https://example.com/1",
		Title:     "Bookmark 1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	bookmark2 := &domain.Bookmark{
		ID:        "id-2",
		URL:       "https://example.com/2",
		Title:     "Bookmark 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name         string
		prePopulate  []*domain.Bookmark
		idToGet      string
		wantErr      error
		wantBookmark *domain.Bookmark
	}{
		{
			name:         "Successfully Get Existing Bookmark",
			prePopulate:  []*domain.Bookmark{bookmark1, bookmark2},
			idToGet:      "id-1",
			wantErr:      nil,
			wantBookmark: bookmark1,
		},
		{
			name:         "Get Non-Existent Bookmark",
			prePopulate:  []*domain.Bookmark{bookmark1},
			idToGet:      "non-existent-id",
			wantErr:      domain.ErrBookmarkNotFound,
			wantBookmark: nil,
		},
		{
			name:         "Get from Empty Repository",
			prePopulate:  []*domain.Bookmark{},
			idToGet:      "any-id",
			wantErr:      domain.ErrBookmarkNotFound,
			wantBookmark: nil,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable for parallel tests
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := newTestRepo()
			ctx := context.Background()

			for _, b := range tt.prePopulate {
				repo.Create(ctx, b)
			}

			gotBookmark, err := repo.GetByID(ctx, tt.idToGet)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && gotBookmark != tt.wantBookmark {
				t.Errorf("GetByID() gotBookmark = %v, want %v", gotBookmark, tt.wantBookmark)
			}
		})
	}
}

func TestInMemoryBookmarkRepository_GetAll(t *testing.T) {
	t.Parallel()

	bookmark1 := &domain.Bookmark{ID: "id-1", URL: "url1", Title: "title1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	bookmark2 := &domain.Bookmark{ID: "id-2", URL: "url2", Title: "title2", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	bookmark3 := &domain.Bookmark{ID: "id-3", URL: "url3", Title: "title3", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	tests := []struct {
		name        string
		prePopulate []*domain.Bookmark
		wantCount   int
	}{
		{
			name:        "Get All from Empty Repository",
			prePopulate: []*domain.Bookmark{},
			wantCount:   0,
		},
		{
			name:        "Get All with One Bookmark",
			prePopulate: []*domain.Bookmark{bookmark1},
			wantCount:   1,
		},
		{
			name:        "Get All with Multiple Bookmarks",
			prePopulate: []*domain.Bookmark{bookmark1, bookmark2, bookmark3},
			wantCount:   3,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable for parallel tests
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := newTestRepo()
			ctx := context.Background()

			for _, b := range tt.prePopulate {
				repo.Create(ctx, b)
			}

			gotBookmarks, err := repo.GetAll(ctx)
			if err != nil {
				t.Errorf("GetAll() unexpected error = %v", err)
			}

			if len(gotBookmarks) != tt.wantCount {
				t.Errorf("GetAll() got %d bookmarks, want %d", len(gotBookmarks), tt.wantCount)
			}

			// Further checks can be added to ensure the correct bookmarks are returned,
			// though for an in-memory map, order is not guaranteed.
			// A simple check could be to build a map of IDs from gotBookmarks and compare with prePopulate IDs.
		})
	}
}

func TestInMemoryBookmarkRepository_Delete(t *testing.T) {
	t.Parallel()

	bookmark1 := &domain.Bookmark{
		ID:        "id-1",
		URL:       "https://example.com/1",
		Title:     "Bookmark 1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	bookmark2 := &domain.Bookmark{
		ID:        "id-2",
		URL:       "https://example.com/2",
		Title:     "Bookmark 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		prePopulate   []*domain.Bookmark
		idToDelete    string
		wantErr       error
		wantRemaining int // Number of bookmarks expected after deletion
	}{
		{
			name:          "Successfully Delete Existing Bookmark",
			prePopulate:   []*domain.Bookmark{bookmark1, bookmark2},
			idToDelete:    "id-1",
			wantErr:       nil,
			wantRemaining: 1,
		},
		{
			name:          "Delete Non-Existent Bookmark",
			prePopulate:   []*domain.Bookmark{bookmark1},
			idToDelete:    "non-existent-id",
			wantErr:       domain.ErrBookmarkNotFound,
			wantRemaining: 1, // Bookmark1 should still be there
		},
		{
			name:          "Delete from Empty Repository",
			prePopulate:   []*domain.Bookmark{},
			idToDelete:    "any-id",
			wantErr:       domain.ErrBookmarkNotFound,
			wantRemaining: 0,
		},
		{
			name:          "Delete last remaining bookmark",
			prePopulate:   []*domain.Bookmark{bookmark1},
			idToDelete:    "id-1",
			wantErr:       nil,
			wantRemaining: 0,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable for parallel tests
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := newTestRepo()
			ctx := context.Background()

			for _, b := range tt.prePopulate {
				repo.Create(ctx, b)
			}

			err := repo.Delete(ctx, tt.idToDelete)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				// Verify the bookmark is actually deleted
				_, getErr := repo.GetByID(ctx, tt.idToDelete)
				if !errors.Is(getErr, domain.ErrBookmarkNotFound) {
					t.Errorf("Delete() expected bookmark to be not found after deletion, got %v", getErr)
				}
			}

			// Verify remaining count
			allBookmarks, _ := repo.GetAll(ctx)
			if len(allBookmarks) != tt.wantRemaining {
				t.Errorf("Delete() got %d remaining bookmarks, want %d", len(allBookmarks), tt.wantRemaining)
			}
		})
	}
}
