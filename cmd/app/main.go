package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/etsrc/goprod/internal/domain"
	"github.com/etsrc/goprod/internal/infra/persistence"
	"github.com/etsrc/goprod/internal/infra/transport/rest"
	"github.com/etsrc/goprod/internal/infra/transport/rest/gen"
	"github.com/etsrc/goprod/internal/service"
)

func main() {
	ctx := context.Background()
	bookmarkRepo := persistence.NewInMemoryBookmarkRepository()
	bookmarkService := service.NewBookmarkService(bookmarkRepo)

	// --- Example Usage ---
	fmt.Println("Setting up Bookmark Service with in-memory repository...")

	// Create a sample bookmark using the factory function
	newBookmark := domain.NewBookmark(
		"https://example.com",
		"Example Domain",
		"A sample bookmark for demonstration",
		[]string{"sample", "test"},
	)

	fmt.Println("Creating a new bookmark...")
	err := bookmarkService.Create(ctx, newBookmark)
	if err != nil {
		log.Fatalf("Failed to create bookmark: %v", err)
	}
	fmt.Printf("Bookmark created successfully with ID: %s\n", newBookmark.ID)

	// Get the bookmark by ID
	fmt.Printf("Fetching bookmark with ID: %s...\n", newBookmark.ID)
	fetchedBookmark, err := bookmarkService.GetByID(ctx, newBookmark.ID)
	if err != nil {
		log.Fatalf("Failed to get bookmark by ID: %v", err)
	}
	fmt.Printf("Fetched bookmark: %+v\n", fetchedBookmark)

	// List all bookmarks
	fmt.Println("Fetching all bookmarks...")
	allBookmarks, err := bookmarkService.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list bookmarks: %v", err)
	}
	fmt.Printf("Found %d bookmark(s).\n", len(allBookmarks))
	for i, b := range allBookmarks {
		fmt.Printf("  %d: ID=%s, Title=%s\n", i+1, b.ID, b.Title)
	}

	// // Delete the bookmark
	// fmt.Printf("Deleting bookmark with ID: %s...\n", newBookmark.ID)
	// err = bookmarkService.Delete(ctx, newBookmark.ID)
	// if err != nil {
	// 	log.Fatalf("Failed to delete bookmark: %v", err)
	// }
	// fmt.Println("Bookmark deleted successfully.")

	// // Verify deletion
	// fmt.Printf("Verifying deletion by fetching bookmark with ID: %s...\n", newBookmark.ID)
	// _, err = bookmarkService.GetByID(ctx, newBookmark.ID)
	// if err != nil {
	// 	if err == domain.ErrBookmarkNotFound {
	// 		fmt.Println("Bookmark not found after deletion, as expected.")
	// 	} else {
	// 		log.Fatalf("Unexpected error when verifying deletion: %v", err)
	// 	}
	// } else {
	// 	log.Fatalf("Bookmark still found after deletion, which is unexpected.")
	// }

	// fmt.Println("In-memory bookmark setup and basic operations successful.")

	// Start Server
	handler := rest.NewBookmarkHandler(bookmarkService)

	mux := http.NewServeMux()
	gen.HandlerFromMux(handler, mux)

	// Just use the mux directly without wrapping it
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("🚀 Server starting on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
