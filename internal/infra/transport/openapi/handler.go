package openapi

import (
	"encoding/json"
	"net/http"

	"github.com/etsrc/goprod/internal/domain"
	"github.com/etsrc/goprod/internal/infra/transport/openapi/gen"
	"github.com/etsrc/goprod/internal/service"
)

// BookmarkHandler implements gen.ServerInterface
type BookmarkHandler struct {
	svc *service.BookmarkService
}

func NewBookmarkHandler(svc *service.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{svc: svc}
}

// GetAllBookmarks handles GET /bookmarks
func (h *BookmarkHandler) GetAllBookmarks(w http.ResponseWriter, r *http.Request) {
	// Assuming your service List method takes context and an optional search string
	bookmarks, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookmarks)
}

// CreateBookmark handles POST /bookmarks
// It uses gen.BookmarkInput as defined in your spec
func (h *BookmarkHandler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	var input gen.BookmarkInput // Use the Input model from generated code
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Map API Input -> Domain Model
	bm := &domain.Bookmark{
		Title: input.Title,
		URL:   input.Url,
	}

	if err := h.svc.Create(r.Context(), bm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bm)
}

// GetBookmarkByID handles GET /bookmarks/{id}
func (h *BookmarkHandler) GetBookmarkByID(w http.ResponseWriter, r *http.Request, id string) {
	bm, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Bookmark not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bm)
}

// DeleteBookmark handles DELETE /bookmarks/{id}
func (h *BookmarkHandler) DeleteBookmark(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
