package rest

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/etsrc/goprod/internal/domain"
	"github.com/etsrc/goprod/internal/mocks"
	"github.com/stretchr/testify/mock"
)

func TestBookmarkHandler_GetAllBookmarks(t *testing.T) {
	t.Parallel()

	bookmarks := []*domain.Bookmark{
		{ID: "1", Title: "Google", URL: "https://google.com"},
		{ID: "2", Title: "Example", URL: "https://example.com"},
	}

	tests := []struct {
		name         string
		mockBehavior func(m *mocks.BookmarkService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "Success",
			mockBehavior: func(m *mocks.BookmarkService) {
				m.On("List", mock.Anything).Return(bookmarks, nil).Once()
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"id":"1","url":"https://google.com","title":"Google","description":"","tags":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"id":"2","url":"https://example.com","title":"Example","description":"","tags":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]` + "\n",
		},
		{
			name: "Service Error",
			mockBehavior: func(m *mocks.BookmarkService) {
				m.On("List", mock.Anything).Return(nil, domain.ErrBookmarkNotFound).Once()
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "bookmark not found\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockSvc := mocks.NewBookmarkService(t)
			tt.mockBehavior(mockSvc)

			handler := NewBookmarkHandler(mockSvc)
			req := httptest.NewRequest("GET", "/bookmarks", nil)
			w := httptest.NewRecorder()

			handler.GetAllBookmarks(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("GetAllBookmarks() status code = %v, want %v", w.Code, tt.expectedCode)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("GetAllBookmarks() body = %q, want %q", w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestBookmarkHandler_CreateBookmark(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		requestBody  string
		mockBehavior func(m *mocks.BookmarkService)
		expectedCode int
		expectedBody string
	}{
		{
			name:        "Success",
			requestBody: `{"title": "New Site", "url": "https://newsite.com"}`,
			mockBehavior: func(m *mocks.BookmarkService) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(b *domain.Bookmark) bool {
					return b.Title == "New Site" && b.URL == "https://newsite.com"
				})).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(1).(*domain.Bookmark)
					arg.ID = "3"
				}).Once()
			},
			expectedCode: http.StatusCreated,
			expectedBody: `{"id":"3","url":"https://newsite.com","title":"New Site","description":"","tags":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}` + "\n",
		},
		{
			name:        "Invalid Request Body",
			requestBody: `{"title": "New Site",`,
			mockBehavior: func(_ *mocks.BookmarkService) {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request body\n",
		},
		{
			name:        "Service Error",
			requestBody: `{"title": "New Site", "url": "https://newsite.com"}`,
			mockBehavior: func(m *mocks.BookmarkService) {
				m.On("Create", mock.Anything, mock.Anything).Return(domain.ErrInvalidURL).Once()
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "the provided URL is invalid\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockSvc := mocks.NewBookmarkService(t)
			tt.mockBehavior(mockSvc)

			handler := NewBookmarkHandler(mockSvc)
			req := httptest.NewRequest("POST", "/bookmarks", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateBookmark(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("CreateBookmark() status code = %v, want %v", w.Code, tt.expectedCode)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("CreateBookmark() body = %q, want %q", w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestBookmarkHandler_GetBookmarkByID(t *testing.T) {
	t.Parallel()

	bookmark := &domain.Bookmark{ID: "1", Title: "Google", URL: "https://google.com"}

	tests := []struct {
		name         string
		bookmarkID   string
		mockBehavior func(m *mocks.BookmarkService)
		expectedCode int
		expectedBody string
	}{
		{
			name:       "Success",
			bookmarkID: "1",
			mockBehavior: func(m *mocks.BookmarkService) {
				m.On("GetByID", mock.Anything, "1").Return(bookmark, nil).Once()
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"1","url":"https://google.com","title":"Google","description":"","tags":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}` + "\n",
		},
		{
			name:       "Not Found",
			bookmarkID: "99",
			mockBehavior: func(m *mocks.BookmarkService) {
				m.On("GetByID", mock.Anything, "99").Return(nil, domain.ErrBookmarkNotFound).Once()
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "Bookmark not found\n",
		},
		{
			name:       "Service Error",
			bookmarkID: "1",
			mockBehavior: func(m *mocks.BookmarkService) {
				m.On("GetByID", mock.Anything, "1").Return(nil, errors.New("Internal server error")).Once()
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "Bookmark not found\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockSvc := mocks.NewBookmarkService(t)
			tt.mockBehavior(mockSvc)

			handler := NewBookmarkHandler(mockSvc)
			req := httptest.NewRequest("GET", "/bookmarks/"+tt.bookmarkID, nil)
			w := httptest.NewRecorder()

			handler.GetBookmarkByID(w, req, tt.bookmarkID)

			if w.Code != tt.expectedCode {
				t.Errorf("GetBookmarkByID() status code = %v, want %v", w.Code, tt.expectedCode)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("GetBookmarkByID() body = %q, want %q", w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestBookmarkHandler_DeleteBookmark(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		bookmarkID   string
		mockBehavior func(m *mocks.BookmarkService)
		expectedCode int
	}{
		{
			name:       "Success",
			bookmarkID: "1",
			mockBehavior: func(m *mocks.BookmarkService) {
				m.On("Delete", mock.Anything, "1").Return(nil).Once()
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:       "Service Error",
			bookmarkID: "1",
			mockBehavior: func(m *mocks.BookmarkService) {
				m.On("Delete", mock.Anything, "1").Return(errors.New("Delete failed")).Once()
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockSvc := mocks.NewBookmarkService(t)
			tt.mockBehavior(mockSvc)

			handler := NewBookmarkHandler(mockSvc)
			req := httptest.NewRequest("DELETE", "/bookmarks/"+tt.bookmarkID, nil)
			w := httptest.NewRecorder()

			handler.DeleteBookmark(w, req, tt.bookmarkID)

			if w.Code != tt.expectedCode {
				t.Errorf("DeleteBookmark() status code = %v, want %v", w.Code, tt.expectedCode)
			}
		})
	}
}
