package main

import (
	"github.com/gmax79/bfservice/internal/buckets"
	"github.com/gmax79/bfservice/internal/netsupport"
	"github.com/gmax79/bfservice/internal/storage"
)

// filter - main objects to filtering login attempts
type filter struct {
	whitelist *netsupport.SubnetsList
	blacklist *netsupport.SubnetsList
	counter   *buckets.AttemptsCounter
	limits    buckets.RatesLimits
	stor      storage.Provider
}

// createFilter - create instance of filter
func createFilter(config RatesAndHostConfig) (*filter, error) {
	stor, err := storage.ConnectRedis(config.RedisHost, config.RedisPassword, config.RedisDB)
	if err != nil {
		return nil, err
	}
	wlProvider, err := stor.CreateSet("abfWhitelist")
	if err != nil {
		return nil, err
	}
	blProvider, err := stor.CreateSet("abdBlacklist")
	if err != nil {
		return nil, err
	}
	var f filter
	f.stor = stor
	f.whitelist, err = netsupport.CreateSubnetsList(wlProvider)
	if err != nil {
		return nil, err
	}
	f.blacklist, err = netsupport.CreateSubnetsList(blProvider)
	if err != nil {
		return nil, err
	}
	f.limits.Login = config.LoginRate
	f.limits.Password = config.PasswordRate
	f.limits.Host = config.IPRate
	f.counter = buckets.CreateCounter(f.limits)
	return &f, nil
}

func (f *filter) CheckLogin(login, password, hostip string) (bool, string, error) {
	var host netsupport.IPAddr
	if err := host.Parse(hostip); err != nil {
		return false, "", err
	}
	if f.blacklist.Check(host) {
		return false, "blocked by blacklist", nil
	}
	if f.whitelist.Check(host) {
		return true, "passed by whitelist", nil
	}
	return f.counter.CheckAndCount(login, password, hostip)
}

func (f *filter) ResetLogin(login, hostip string) (bool, error) {
	return f.counter.Reset(login, hostip)
}

func (f *filter) AddWhiteList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	added, err := f.whitelist.Add(snet)
	return added, err
}

func (f *filter) DeleteWhiteList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	deleted, err := f.whitelist.Delete(snet)
	return deleted, err
}

func (f *filter) AddBlackList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	added, err := f.blacklist.Add(snet)
	return added, err
}

func (f *filter) DeleteBlackList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	deleted, err := f.blacklist.Delete(snet)
	return deleted, err
}

func (f *filter) GetLimits() buckets.RatesLimits {
	return f.limits
}

func (f *filter) Close() error {
	return f.stor.Close()
}
