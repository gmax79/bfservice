package buckets

import "github.com/gmax79/bfservice/internal/hostslist"

type filterImpl struct {
	whitelist *hostslist.Filter
	blacklist *hostslist.Filter
}

// CreateFilter - create instance of filter
func CreateFilter() Filter {
	f := filterImpl{}
	f.whitelist = hostslist.CreateFilter()
	f.blacklist = hostslist.CreateFilter()
	return &f
}

func (f *filterImpl) CheckLogin(login, password, ip string) error {
	return nil
}

func (f *filterImpl) ResetLogin(login, ip string) error {
	return nil
}

func (f *filterImpl) AddWhiteList(ipsubnet string) error {
	return nil
}

func (f *filterImpl) DeleteWhiteList(ipsubnet string) error {
	return nil
}

func (f *filterImpl) AddBlackList(ipsubnet string) error {
	return nil
}

func (f *filterImpl) DeleteBlackList(ipsubnet string) error {
	return nil
}
