package in_memory

import (
	"fmt"
	"strings"

	"sync"

	"github.com/rs/zerolog/log"
)

const (
	errForbiddenWords = "forbidden word"
	errWordNotFound   = "word not found"

	limitWords     = 3
	logCleanGC     = "cleaning the garbage collector"
	logStorage     = "current storage: %v"
	logWordAddINFO = "word %s succesfully inserted to storage"
	logWordGetINFO = "most frequent word with this prefix '%s' is %s"
)

// forbiddenWords serves as an example list of words that are disallowed
// for insertion. Given that our in-memory storage simply adds words to a map,
// there's a minimal chance of operational errors. However, to adhere to our
// interface which expects an error return, we use this list to simulate potential
// error scenarios during the Insert() operation.
var forbiddenWords = []string{"bad", "words", "example"}

type InMemoryStorage struct {
	mu sync.RWMutex

	WordsStorage map[string]int
}

func NewInMemoryStorage() *InMemoryStorage {
	storage := &InMemoryStorage{
		WordsStorage: make(map[string]int, 0),
	}

	return storage
}

// TODO
func (s *InMemoryStorage) Insert(word string) error {
	// Simulated error for this challenge
	for _, w := range forbiddenWords {
		if word == w {
			return fmt.Errorf(errForbiddenWords)
		}
	}

	s.mu.Lock()
	s.WordsStorage[word]++
	s.mu.Unlock()

	log.Info().Msgf(logWordAddINFO, word)

	return nil
}

// TODO
func (s *InMemoryStorage) FindFrequentByPrefix(prefix string) (string, error) {
	var maxWord string
	var maxCount int

	s.mu.RLock()
	for word, count := range s.WordsStorage {
		if strings.HasPrefix(word, prefix) && count > maxCount {
			maxWord = word
			maxCount = count
		}
	}
	s.mu.RUnlock()

	if strings.TrimSpace(maxWord) == "" {
		return "", fmt.Errorf(errWordNotFound)
	}

	log.Info().Msgf(logWordGetINFO, prefix, maxWord)

	return maxWord, nil
}

// CleanGarbageCollector trims infrequent words (e.g., count of 3 for testing challenge, but 500 in real scenarios)
// from InMemoryStore to save memory and enhance performance.
// Especially useful when storage swells, limiting iteration overhead over vast entries like 10,000 words.
func (s *InMemoryStorage) CleanGarbageCollector() {
	log.Info().Msg(logCleanGC)
	log.Info().Msgf(logStorage, s.WordsStorage)

	s.mu.Lock()
	defer s.mu.Unlock()

	for word := range s.WordsStorage {
		if s.WordsStorage[word] <= limitWords {
			delete(s.WordsStorage, word)
		}
	}
}
