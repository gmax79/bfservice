package storage

// InMemoryStorage - create to in memory storage
func InMemoryStorage() (Provider, error) {
	var p memoryProvider
	return &p, nil
}

type memoryProvider struct {
}

func (p *memoryProvider) CreateSet(id string) (SetProvider, error) {
	return createMemorySetProvider(), nil
}

func (p *memoryProvider) Close() error {
	return nil
}
