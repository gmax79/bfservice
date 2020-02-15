package storage

// Provider - load/save settings in extenal storage
type Provider interface {
	CreateSet(id string) (SetProvider, error)
	Close() error
}

// StringIterator - interface to iterate over list
type StringIterator func() (string, bool) // string, end of array flag

// SetProvider - interface to working with set in external storage
type SetProvider interface {
	Add(item string) (bool, error)     // added flag
	Delete(item string) (bool, error)  // deleted flag
	Load(f func(v string) error) error // load set from storage via functor
}
