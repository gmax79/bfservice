package storage

// StringIterator - interface to iterate over list
type StringIterator func() (string, bool)

// ListProvider - interface to working with list in external storage
type ListProvider interface {
	Add(item string) (bool, error)
	Delete(item string) (bool, error)
	Iterator() (StringIterator, error)
}
