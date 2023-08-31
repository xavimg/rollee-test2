package storage

// Storer is anything that can store words into a storage.
type WordStorager interface {
	Insert(word string) error
	FindFrequentByPrefix(prefix string) (string, error)
}
