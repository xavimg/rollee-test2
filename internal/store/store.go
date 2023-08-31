package store

// Storer is anything that can store words into a storage.
type WordStorer interface {
	Insert(word string) error
	FindFrequentByPrefix(prefix string) (string, error)
}
