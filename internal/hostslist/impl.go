package hostslist

// Filter - object, to filtering hosts by some rules
type Filter struct {
	list map[subnet]struct{}
}

// CreateFilter - creates hosts filter object
func CreateFilter() *Filter {
	f := &Filter{
		list: make(map[subnet]struct{}),
	}
	return f
}

// Add - add subnet into list
func (f *Filter) Add(s subnet) bool {
	_, exist := f.list[s]
	f.list[s] = struct{}{}
	return exist
}

// Delete - delete subnet from list
func (f *Filter) Delete(s subnet) bool {
	_, exist := f.list[s]
	delete(f.list, s)
	return exist
}

// Check - check host in list
func (f *Filter) Check(host ip) bool {
	for s := range f.list {
		if s.Check(host) {
			return true
		}
	}
	return false
}
