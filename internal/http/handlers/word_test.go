package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"words/internal/service/word"
	"words/internal/storage/in_memory"

	"github.com/go-chi/chi"
)

func TestAddWord(t *testing.T) {
	repository := in_memory.NewInMemoryStorage()
	mockService := &word.WordService{
		Repository: repository,
	}
	handler := NewHandler(mockService)
	r := chi.NewRouter()

	r.Post("/words/{word}", handler.AddWord)

	// Test successful add for the first time
	req, _ := http.NewRequest("POST", "/words/testword", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	// Test successful add for the second time (duplicates allowed)
	req, _ = http.NewRequest("POST", "/words/testword", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestFrequentWordByPrefix(t *testing.T) {
	repository := in_memory.NewInMemoryStorage()
	mockService := &word.WordService{
		Repository: repository,
	}
	handler := NewHandler(mockService)
	r := chi.NewRouter()

	// First, let's add a word "testword" to the repository
	err := mockService.AddWord("testword")
	if err != nil {
		t.Fatalf("Failed to add word to mock service: %v", err)
	}

	r.Get("/words/{prefix}", handler.FrequentWordByPrefix)

	// Test successful find
	req, _ := http.NewRequest("GET", "/words/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "testword" {
		t.Errorf("expected word '%s', got '%s'", "testword", rec.Body.String())
	}

	// Test error on find for a non-existing prefix
	req, _ = http.NewRequest("GET", "/words/nonexistentprefix", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestValidateWordFormat(t *testing.T) {
	tests := []struct {
		name      string
		word      string
		want      bool
		wantError bool
	}{
		{
			name:      "valid word",
			word:      "hello",
			want:      true,
			wantError: false,
		},
		{
			name:      "invalid word with numbers",
			word:      "hello123",
			want:      false,
			wantError: true,
		},
		{
			name:      "invalid word with special characters",
			word:      "hello@world",
			want:      false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateWordFormat(tt.word, regexPattern)
			if (err != nil) != tt.wantError {
				t.Errorf("validateWordFormat() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if got != tt.want {
				t.Errorf("validateWordFormat() got = %v, want %v", got, tt.want)
			}
		})
	}
}
