package storage

// StringIterator - interface to iterate over list
type StringIterator func() (string, bool)

// SetProvider - interface to working with set in external storage
type SetProvider interface {
	Add(item string) (bool, error)
	Delete(item string) (bool, error)
	Iterator() (StringIterator, error)
}
