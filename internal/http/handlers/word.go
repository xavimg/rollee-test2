package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"words/internal/service/word"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

const (
	errBadRequest      = "bad request"
	errInvalidInput    = "invalid format input received"
	errWordNotResolved = "word not found"

	logAddWordERROR = "problem adding new word: %s. error: %v"
	logGetWordERROR = "error retrieving word for prefix '%s'. error: %v"

	prefixParam  = "prefix"
	regexPattern = `^[a-zA-Z]+$`
	wordParam    = "word"
)

type WordHandler struct {
	wordService *word.WordService
}

func NewHandler(s *word.WordService) *WordHandler {
	return &WordHandler{wordService: s}
}

func (h *WordHandler) AddWord(w http.ResponseWriter, r *http.Request) {
	word := chi.URLParam(r, wordParam)

	if ok, err := validateWordFormat(word, regexPattern); !ok {
		log.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.wordService.AddWord(word); err != nil {
		log.Error().Msgf(logAddWordERROR, word, err)
		http.Error(w, errBadRequest, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *WordHandler) FrequentWordByPrefix(w http.ResponseWriter, r *http.Request) {
	prefix := chi.URLParam(r, prefixParam)

	if ok, err := validateWordFormat(prefix, regexPattern); !ok {
		log.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	word, err := h.wordService.GetMostFrequentByPrefix(prefix)
	if err != nil {
		log.Error().Msgf(logGetWordERROR, prefix, err)
		http.Error(w, errWordNotResolved, http.StatusNotFound)
		return
	}

	w.Write([]byte(word))
}

// validateWordFormat checks if the given word matches a predefined regex pattern.
func validateWordFormat(word, regexPattern string) (bool, error) {
	matched, err := regexp.MatchString(regexPattern, word)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, errors.New(errInvalidInput)
	}

	return true, nil
}
