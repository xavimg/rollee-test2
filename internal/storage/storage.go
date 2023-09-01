package storage

// WordRepository is anything that can store words into a storage.
type WordRepository interface {
	Insert(word string) error
	FindFrequentByPrefix(prefix string) (string, error)
}
