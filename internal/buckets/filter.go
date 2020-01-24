package buckets

// Filter - main objects to filtering login attempts
type Filter interface {
	CheckLogin(login, password, hostip string) error
	ResetLogin(login, hostip string) error
	AddWhiteList(subnetip string) error
	DeleteWhiteList(subnetip string) error
	AddBlackList(subnetip string) error
	DeleteBlackList(subnetip string) error
}
