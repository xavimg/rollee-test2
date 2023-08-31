package storage

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	logWordAddINFO = "word %s succesfully inserted to storage"
)

// forbiddenWords serves as an example list of words that are disallowed
// for insertion. Given that our in-memory storage simply adds words to a map,
// there's a minimal chance of operational errors. However, to adhere to our
// interface which expects an error return, we use this list to simulate potential
// error scenarios during the Insert() operation.
var forbiddenWords = []string{"bad", "words", "example"}

type InMemoryStore struct {
	WordsStore map[string]int

	mux        sync.RWMutex
	GcInterval time.Duration
}

func NewInMemoryStore(GcInterval time.Duration) *InMemoryStore {
	store := &InMemoryStore{
		WordsStore: make(map[string]int, 0),
		GcInterval: GcInterval,
	}

	go store.cleanGarbageCollector()

	return store
}

func (s *InMemoryStore) Insert(word string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Simulated error for this challenge
	for _, w := range forbiddenWords {
		if word == w {
			return fmt.Errorf("forbidden word")
		}
	}

	s.WordsStore[word]++
	log.Info().Msgf(logWordAddINFO, word)
	return nil
}

func (s *InMemoryStore) FindFrequentByPrefix(prefix string) (string, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	var maxWord string
	var maxCount int
	for word, count := range s.WordsStore {
		if strings.HasPrefix(word, prefix) && count > maxCount {
			maxWord = word
			maxCount = count
		}
	}
	if maxWord == "" {
		return "", fmt.Errorf("word notfound")
	}

	log.Info().Msgf("the most frequent word with this prefix '%s' is %s", prefix, maxWord)
	return maxWord, nil
}

// cleanGarbageCollector periodically removes words from the InMemoryStore
// that have a count of 1. This serves as a cleanup mechanism to free up memory
// by removing less frequently used words. The cleanup interval is determined by
// the `GcInterval` field of the `InMemoryStore` struct. This method should typically
// be run as a goroutine since it will loop indefinitely, cleaning up words at the
// specified interval. Also, we dont want to go trought so many words when we iterate
// the map, so this helps to performance too.
func (s *InMemoryStore) cleanGarbageCollector() {
	ticker := time.NewTicker(s.GcInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Info().Msg("Cleaning the garbage collector")
			log.Info().Msgf("Actual storage capacity: %v", s.WordsStore)
			s.mux.Lock()
			for word := range s.WordsStore {
				if s.WordsStore[word] == 1 {
					delete(s.WordsStore, word)
				}
			}
			s.mux.Unlock()
		}
	}
}
