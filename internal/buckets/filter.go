package buckets

// Filter - main objects to filtering login attempts
type Filter interface {
	CheckLogin(login, password, ip string) error
	ResetLogin(login, ip string) error
	AddWhiteList(ipsubnet string) error
	DeleteWhiteList(ipsubnet string) error
	AddBlackList(ipsubnet string) error
	DeleteBlackList(ipsubnet string) error
}

// CreateFilter - create instance of filter
func CreateFilter() Filter {
	return &filterImpl{}
}
