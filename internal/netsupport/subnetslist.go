package netsupport

import "github.com/gmax79/bfservice/internal/storage"

// SubnetsList - object, to filtering hosts by some rules
type SubnetsList struct {
	list map[Subnet]struct{}
	prov storage.SetProvider
}

// CreateSubnetsList - creates hosts/subnets list
func CreateSubnetsList(prov storage.SetProvider) (*SubnetsList, error) {
	f := &SubnetsList{
		list: make(map[Subnet]struct{}),
		prov: prov,
	}
	// load table in memory
	iter, err := prov.Iterator()
	if err != nil {
		return nil, err
	}
	s, ok := iter()
	for ; ok; s, ok = iter() {
		var snet Subnet
		err := snet.Parse(s)
		if err != nil {
			return nil, err
		}
		f.list[snet] = struct{}{}
	}
	return f, nil
}

// Add - add subnet into list
func (s *SubnetsList) Add(snet Subnet) (bool, error) {
	_, exist := s.list[snet]
	if exist {
		return true, nil
	}
	s.list[snet] = struct{}{}
	return s.prov.Add(snet.String())
}

// Delete - delete subnet from list
func (s *SubnetsList) Delete(snet Subnet) (bool, error) {
	_, exist := s.list[snet]
	if !exist {
		return false, nil
	}
	delete(s.list, snet)
	return s.prov.Delete(snet.String())
}

// Check - check host in list
func (s *SubnetsList) Check(host IPAddr) bool {
	for snet := range s.list {
		if snet.Check(host) {
			return true
		}
	}
	return false
}
