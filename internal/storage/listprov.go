package storage

// ListProvider - interface to working with list in external storage
type ListProvider interface {
	Add(item string) (bool, error)
	Delete(item string) (bool, error)
	Count() (int, error)
	Get(index int) (string, error)
}
