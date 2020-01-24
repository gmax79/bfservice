package netsupport

// SubnetsList - object, to filtering hosts by some rules
type SubnetsList struct {
	list map[Subnet]struct{}
}

// CreateSubnetsList - creates hosts/subnets list
func CreateSubnetsList() *SubnetsList {
	f := &SubnetsList{
		list: make(map[Subnet]struct{}),
	}
	return f
}

// Add - add subnet into list
func (s *SubnetsList) Add(snet Subnet) {
	s.list[snet] = struct{}{}
}

// Delete - delete subnet from list
func (s *SubnetsList) Delete(snet Subnet) {
	delete(s.list, snet)
}

// Exist - check subnet existence in list
func (s *SubnetsList) Exist(snet Subnet) bool {
	_, exist := s.list[snet]
	return exist
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
