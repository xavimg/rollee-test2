package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

type MockWordService struct {
	addWordErr             error
	getMostFrequentByError error
	frequentWord           string
}

func (m *MockWordService) AddWord(word string) error {
	return m.addWordErr
}

func (m *MockWordService) GetMostFrequentByPrefix(prefix string) (string, error) {
	if m.getMostFrequentByError != nil {
		return "", m.getMostFrequentByError
	}
	return m.frequentWord, nil
}

func TestAddWord(t *testing.T) {
	mockService := &MockWordService{}
	handler := NewHandler(mockService)
	r := chi.NewRouter()

	r.Post("/words/{word}", handler.AddWord)

	// Test successful add
	req, _ := http.NewRequest("POST", "/words/testword", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	// Test error on add
	expectedError := errors.New("mock add error")
	mockService.addWordErr = expectedError
	req, _ = http.NewRequest("POST", "/words/testword", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestFrequentWordByPrefix(t *testing.T) {
	mockService := &MockWordService{
		frequentWord: "testword",
	}
	handler := NewHandler(mockService)
	r := chi.NewRouter()

	r.Get("/words/{prefix}", handler.FrequentWordByPrefix)

	// Test successful find
	req, _ := http.NewRequest("GET", "/words/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "testword" {
		t.Errorf("expected word 'testword', got %s", rec.Body.String())
	}

	// Test error on find
	expectedError := errors.New("mock find error")
	mockService.getMostFrequentByError = expectedError
	req, _ = http.NewRequest("GET", "/words/test", nil)
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
