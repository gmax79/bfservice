package storage

func createMemorySetProvider() *memorySetProvider {
	var p memorySetProvider
	p.set = make(map[string]struct{})
	return &p
}

type memorySetProvider struct {
	set map[string]struct{}
}

func (p *memorySetProvider) Add(item string) (bool, error) {
	_, ok := p.set[item]
	if ok {
		return false, nil
	}
	p.set[item] = struct{}{}
	return true, nil
}

func (p *memorySetProvider) Delete(item string) (bool, error) {
	_, ok := p.set[item]
	if ok {
		delete(p.set, item)
		return true, nil
	}
	return false, nil
}

func (p *memorySetProvider) Load(f func(v string) error) error {
	// memory set not load data from somethere
	return nil
}
