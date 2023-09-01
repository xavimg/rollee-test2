package word

import (
	"strings"

	"words/internal/storage"
)

type WordService struct {
	Repository storage.WordRepository
}

func NewWordService(repo storage.WordRepository) *WordService {
	return &WordService{Repository: repo}
}

func (ws *WordService) AddWord(word string) error {
	word = strings.ToLower(word)
	return ws.Repository.Insert(word)
}

func (ws *WordService) GetMostFrequentByPrefix(prefix string) (string, error) {
	prefix = strings.ToLower(prefix)
	return ws.Repository.FindFrequentByPrefix(prefix)
}
