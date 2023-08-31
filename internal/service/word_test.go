package service

import (
	"errors"
	"testing"
)

// MockWordStorer is a mock implementation of the WordStorer interface.
type MockWordStorer struct {
	insertErr  error
	findErr    error
	frequentBy string
}

func (m *MockWordStorer) Insert(word string) error {
	return m.insertErr
}

func (m *MockWordStorer) FindFrequentByPrefix(prefix string) (string, error) {
	if m.findErr != nil {
		return "", m.findErr
	}
	return m.frequentBy, nil
}

func TestAddWord(t *testing.T) {
	mockStore := &MockWordStorer{}
	service := NewWordService(mockStore)

	// Test successful insert
	if err := service.AddWord("test"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test error on insert
	expectedError := errors.New("mock insert error")
	mockStore.insertErr = expectedError
	if err := service.AddWord("test"); err != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

func TestGetMostFrequentByPrefix(t *testing.T) {
	mockStore := &MockWordStorer{
		frequentBy: "testword",
	}
	service := NewWordService(mockStore)

	// Test successful find
	word, err := service.GetMostFrequentByPrefix("test")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if word != "testword" {
		t.Errorf("expected word 'testword', got %v", word)
	}

	// Test error on find
	expectedError := errors.New("mock find error")
	mockStore.findErr = expectedError
	_, err = service.GetMostFrequentByPrefix("test")
	if err != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}
