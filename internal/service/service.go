package service

// WordServicer is any service that can do this methods.
type WordServicer interface {
	AddWord(string) error
	GetMostFrequentByPrefix(string) (string, error)
}
