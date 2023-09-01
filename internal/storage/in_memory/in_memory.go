package in_memory

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	logWordAddINFO = "word %s succesfully inserted to storage"
	logWordGetINFO = "most frequent word with this prefix '%s' is %s"
	logCleanGC     = "cleaning the garbage collector"
	logStorage     = "current storage: %v"
)

// forbiddenWords serves as an example list of words that are disallowed
// for insertion. Given that our in-memory storage simply adds words to a map,
// there's a minimal chance of operational errors. However, to adhere to our
// interface which expects an error return, we use this list to simulate potential
// error scenarios during the Insert() operation.
var forbiddenWords = []string{"bad", "words", "example"}

type InMemoryStorage struct {
	WordsStorage map[string]int

	mu         sync.RWMutex
	GcInterval time.Duration
}

func NewInMemoryStorage(GcInterval time.Duration) *InMemoryStorage {
	storage := &InMemoryStorage{
		WordsStorage: make(map[string]int, 0),
		GcInterval:   GcInterval,
	}

	// go storage.CleanGarbageCollector()

	return storage
}

func (s *InMemoryStorage) Insert(word string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Simulated error for this challenge
	for _, w := range forbiddenWords {
		if word == w {
			return fmt.Errorf("forbidden word")
		}
	}

	s.WordsStorage[word]++
	log.Info().Msgf(logWordAddINFO, word)
	return nil
}

func (s *InMemoryStorage) FindFrequentByPrefix(prefix string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var maxWord string
	var maxCount int
	for word, count := range s.WordsStorage {
		if strings.HasPrefix(word, prefix) && count > maxCount {
			maxWord = word
			maxCount = count
		}
	}
	if maxWord == "" {
		return "", fmt.Errorf("word notfound")
	}

	log.Info().Msgf(logWordGetINFO, prefix, maxWord)
	return maxWord, nil
}

// cleanGarbageCollector trims infrequent words (e.g., count of 3 for testing challenge, but 500 in real scenarios)
// from InMemoryStore to save memory and enhance performance.
// Especially useful when storage swells, limiting iteration overhead over vast entries like 10,000 words.
func (s *InMemoryStorage) CleanGarbageCollector() {
	ticker := time.NewTicker(s.GcInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Info().Msg(logCleanGC)
			log.Info().Msgf(logStorage, s.WordsStorage)
			s.mu.Lock()
			for word := range s.WordsStorage {
				if s.WordsStorage[word] == 3 {
					delete(s.WordsStorage, word)
				}
			}
			s.mu.Unlock()
		}
	}
}
