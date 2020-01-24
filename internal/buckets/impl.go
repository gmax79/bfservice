package buckets

import "github.com/gmax79/bfservice/internal/netsupport"

type filterImpl struct {
	whitelist *netsupport.SubnetsList
	blacklist *netsupport.SubnetsList
}

// CreateFilter - create instance of filter
func CreateFilter() Filter {
	f := filterImpl{}
	f.whitelist = netsupport.CreateSubnetsList()
	f.blacklist = netsupport.CreateSubnetsList()
	return &f
}

func (f *filterImpl) CheckLogin(login, password, hostip string) error {

	var host netsupport.IPAddr
	if err := host.Parse(hostip); err != nil {
		return err
	}
	if f.blacklist.Check(host) {

	}

	return nil
}

func (f *filterImpl) ResetLogin(login, hostip string) error {
	return nil
}

func (f *filterImpl) AddWhiteList(subnetip string) error {
	return nil
}

func (f *filterImpl) DeleteWhiteList(subnetip string) error {
	return nil
}

func (f *filterImpl) AddBlackList(subnetip string) error {
	return nil
}

func (f *filterImpl) DeleteBlackList(subnetip string) error {
	return nil
}
