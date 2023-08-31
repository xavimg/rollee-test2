package service

import (
	"strings"
	"words/internal/store"
)

type WordService struct {
	wordStore store.WordStorer
}

func NewWordService(s store.WordStorer) *WordService {
	return &WordService{wordStore: s}
}

func (ws *WordService) AddWord(word string) error {
	word = strings.ToLower(word)
	return ws.wordStore.Insert(word)
}

func (ws *WordService) GetMostFrequentByPrefix(prefix string) (string, error) {
	prefix = strings.ToLower(prefix)
	return ws.wordStore.FindFrequentByPrefix(prefix)
}
