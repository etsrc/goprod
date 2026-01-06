package service_test

import (
	"context"
	"testing"

	"github.com/etsrc/goprod/internal/domain"
	persistence "github.com/etsrc/goprod/internal/infra/persistence/inmem"
	"github.com/etsrc/goprod/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookmarkIntegration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := persistence.NewInMemoryBookmarkRepository()
	svc := service.NewBookmarkService(repo)

	t.Run("Full Lifecycle", func(t *testing.T) {
		t.Parallel()
		// 1. Create
		newBookmark := domain.NewBookmark(
			"https://example.com",
			"Example Domain",
			"A sample bookmark",
			[]string{"test"},
		)
		err := svc.Create(ctx, newBookmark)
		require.NoError(t, err)

		// 2. Get
		fetched, err := svc.GetByID(ctx, newBookmark.ID)
		require.NoError(t, err)
		assert.Equal(t, newBookmark.Title, fetched.Title)

		// 3. List
		all, err := svc.List(ctx)
		require.NoError(t, err)
		assert.Len(t, all, 1)

		// 4. Delete
		err = svc.Delete(ctx, newBookmark.ID)
		require.NoError(t, err)

		// 5. Verify Deletion
		_, err = svc.GetByID(ctx, newBookmark.ID)
		assert.ErrorIs(t, err, domain.ErrBookmarkNotFound)
	})
}
