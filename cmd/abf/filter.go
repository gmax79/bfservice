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
	stor      *storage.Provider
}

// createFilter - create instance of filter
func createFilter(config RatesAndHostConfig) (*filter, error) {
	stor, err := storage.ConnectRedis(config.RedisHost, config.RedisPassword, config.RedisDB)
	if err != nil {
		return nil, err
	}
	f := filter{}
	f.whitelist = netsupport.CreateSubnetsList()
	f.blacklist = netsupport.CreateSubnetsList()
	f.limits.Login = config.LoginRate
	f.limits.Password = config.PasswordRate
	f.limits.Host = config.IPRate
	f.counter = buckets.CreateCounter(f.limits)
	f.stor = stor
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
	added := f.whitelist.Add(snet)
	return added, nil
}

func (f *filter) DeleteWhiteList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	if !f.whitelist.Exist(snet) {
		return false, nil
	}
	deleted := f.whitelist.Delete(snet)
	return deleted, nil
}

func (f *filter) AddBlackList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	added := f.blacklist.Add(snet)
	return added, nil
}

func (f *filter) DeleteBlackList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	if !f.blacklist.Exist(snet) {
		return false, nil
	}
	deleted := f.blacklist.Delete(snet)
	return deleted, nil
}

func (f *filter) GetLimits() buckets.RatesLimits {
	return f.limits
}
