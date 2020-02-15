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
	// load from provider into cache
	functor := func(v string) error {
		var snet Subnet
		err := snet.Parse(v)
		if err != nil {
			return err
		}
		f.list[snet] = struct{}{}
		return nil
	}
	err := prov.Load(functor)
	return f, err
}

// Add - add subnet into list, result added status or error
func (s *SubnetsList) Add(snet Subnet) (bool, error) {
	s.list[snet] = struct{}{}
	return s.prov.Add(snet.String())
}

// Delete - delete subnet from list, result deleted status or error
func (s *SubnetsList) Delete(snet Subnet) (bool, error) {
	delete(s.list, snet)
	return s.prov.Delete(snet.String())
}

// Check - check host in list, use only cache
func (s *SubnetsList) Check(host IPAddr) bool {
	for snet := range s.list {
		if snet.Check(host) {
			return true
		}
	}
	return false
}
