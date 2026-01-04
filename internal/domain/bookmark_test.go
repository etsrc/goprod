package domain

import (
	"errors" // To compare error types correctly
	"testing"
	"time"
)

func TestBookmark_Validate(t *testing.T) {
	t.Parallel() // Make the test function itself run in parallel

	// Define a valid bookmark for reuse in tests where only specific fields are changed
	validBookmark := Bookmark{
		ID:          "valid-id",
		URL:         "https://www.example.com",
		Title:       "A Valid Title",
		Description: "This is a valid bookmark.",
		Tags:        []string{"valid", "test"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tests := []struct {
		name     string
		bookmark Bookmark
		wantErr  error
	}{
		{
			name:     "Valid Bookmark",
			bookmark: validBookmark, // Use the pre-defined valid bookmark
			wantErr:  nil,
		},
		{
			name: "Empty Title",
			bookmark: func() Bookmark {
				b := validBookmark
				b.Title = ""
				return b
			}(),
			wantErr: ErrTitleTooShort,
		},
		{
			name: "Title Too Short (1 Char)",
			bookmark: func() Bookmark {
				b := validBookmark
				b.Title = "a"
				return b
			}(),
			wantErr: ErrTitleTooShort,
		},
		{
			name: "Title Too Short (2 Chars)",
			bookmark: func() Bookmark {
				b := validBookmark
				b.Title = "ab"
				return b
			}(),
			wantErr: ErrTitleTooShort,
		},
		{
			name: "Title Exactly 3 Chars",
			bookmark: func() Bookmark {
				b := validBookmark
				b.Title = "abc"
				return b
			}(),
			wantErr: nil,
		},
		{
			name: "Title With Leading/Trailing Spaces Trimming Too Short",
			bookmark: func() Bookmark {
				b := validBookmark
				b.Title = "  ab  " // Trims to "ab"
				return b
			}(),
			wantErr: ErrTitleTooShort,
		},
		{
			name: "Title With Leading/Trailing Spaces Trimming Valid",
			bookmark: func() Bookmark {
				b := validBookmark
				b.Title = "  A Valid Title  " // Trims to "A Valid Title"
				return b
			}(),
			wantErr: nil,
		},
		{
			name: "Invalid URL Malformed",
			bookmark: func() Bookmark {
				b := validBookmark
				b.URL = "invalid-url"
				return b
			}(),
			wantErr: ErrInvalidURL,
		},
		{
			name: "Invalid URL No Scheme",
			bookmark: func() Bookmark {
				b := validBookmark
				b.URL = "www.example.com" // url.ParseRequestURI requires a scheme
				return b
			}(),
			wantErr: ErrInvalidURL,
		},
		{
			name: "Valid URL",
			bookmark: func() Bookmark {
				b := validBookmark
				b.URL = "http://localhost:8080/path?query=value"
				return b
			}(),
			wantErr: nil,
		},
		{
			name: "Valid URL with HTTPS",
			bookmark: func() Bookmark {
				b := validBookmark
				b.URL = "https://example.com/some/path"
				return b
			}(),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Make each test case run in parallel

			err := tt.bookmark.Validate()

			// Check for expected error type
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Bookmark_Validate() unexpected error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err == nil {
					t.Errorf("Bookmark_Validate() error = %v, wantErr %v", err, tt.wantErr)
				} else if !errors.Is(err, tt.wantErr) { // Use errors.Is for checking wrapped errors or specific error types
					t.Errorf("Bookmark_Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
