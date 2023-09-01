package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"words/internal/service"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

const (
	errMatchingRegex   = "error matching regex for word '%s'. Error: %v"
	errInvalidWord     = "invalid word format received: %s"
	errInternalServer  = "internal server error"
	errBadRequest      = "bad request"
	errInvalidInput    = "invalid input"
	errWordNotResolved = "word not found"

	logAddWordERROR = "problem adding new word: %s. error: %v"
	logGetWordERROR = "error retrieving word for prefix '%s'. error: %v"

	regexPattern = `^[a-zA-Z]+$`
)

type WordHandler struct {
	wordService service.WordServicer
}

func NewHandler(s service.WordServicer) *WordHandler {
	return &WordHandler{wordService: s}
}

func (h *WordHandler) AddWord(w http.ResponseWriter, r *http.Request) {
	word := chi.URLParam(r, "word")

	if ok, err := validateWordFormat(word, regexPattern); !ok {
		log.Error().Msg(err.Error())
		if err.Error() == fmt.Sprintf(errInvalidWord, word) {
			http.Error(w, errInvalidInput, http.StatusBadRequest)
			return
		}
	}

	if err := h.wordService.AddWord(word); err != nil {
		log.Error().Msgf(logAddWordERROR, word, err)
		http.Error(w, errBadRequest, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *WordHandler) FrequentWordByPrefix(w http.ResponseWriter, r *http.Request) {
	prefix := chi.URLParam(r, "prefix")

	if ok, err := validateWordFormat(prefix, regexPattern); !ok {
		log.Error().Msg(err.Error())
		if err.Error() == fmt.Sprintf(errInvalidWord, prefix) {
			http.Error(w, errInvalidInput, http.StatusBadRequest)
		}
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
		return false, fmt.Errorf(errMatchingRegex, word, err)
	}
	if !matched {
		return false, fmt.Errorf(errInvalidWord, word)
	}
	return true, nil
}
