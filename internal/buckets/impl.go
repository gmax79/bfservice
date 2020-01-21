package buckets

type filterImpl struct {
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
