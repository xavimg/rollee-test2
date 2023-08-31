package store

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
	wordsStore map[string]int

	mux        sync.RWMutex
	gcInterval time.Duration
}

func NewInMemoryStore(gcInterval time.Duration) *InMemoryStore {
	store := &InMemoryStore{
		wordsStore: make(map[string]int),
		gcInterval: gcInterval,
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

	s.wordsStore[word]++
	log.Info().Msgf(logWordAddINFO, word)
	return nil
}

func (s *InMemoryStore) FindFrequentByPrefix(prefix string) (string, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	var maxWord string
	var maxCount int
	for word, count := range s.wordsStore {
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
// the `gcInterval` field of the `InMemoryStore` struct. This method should typically
// be run as a goroutine since it will loop indefinitely, cleaning up words at the
// specified interval. Also, we dont want to go trought so many words when we iterate
// the map, so this helps to performance too.
func (s *InMemoryStore) cleanGarbageCollector() {
	ticker := time.NewTicker(s.gcInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Info().Msg("Cleaning the garbage collector")
			log.Info().Msgf("Actual storage capacity: %v", s.wordsStore)
			s.mux.Lock()
			for word := range s.wordsStore {
				if s.wordsStore[word] == 1 {
					delete(s.wordsStore, word)
				}
			}
			s.mux.Unlock()
		}
	}
}
