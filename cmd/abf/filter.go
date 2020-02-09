package main

import (
	"sync"

	"github.com/gmax79/bfservice/internal/buckets"
	"github.com/gmax79/bfservice/internal/netsupport"
)

// filter - main objects to filtering login attempts
type filter struct {
	whitelist *netsupport.SubnetsList
	blacklist *netsupport.SubnetsList
	wmx, bmx  *sync.Mutex
	counter   *buckets.AttemptsCounter
	limits    buckets.RatesLimits
}

// createFilter - create instance of filter
func createFilter(config RatesAndHostConfig) *filter {
	f := filter{}
	f.whitelist = netsupport.CreateSubnetsList()
	f.blacklist = netsupport.CreateSubnetsList()
	f.wmx = &sync.Mutex{}
	f.bmx = &sync.Mutex{}
	f.limits.Login = config.LoginRate
	f.limits.Password = config.PasswordRate
	f.limits.Host = config.IPRate
	f.counter = buckets.CreateCounter(f.limits)
	return &f
}

func (f *filter) CheckLogin(login, password, hostip string) (bool, string, error) {
	var host netsupport.IPAddr
	if err := host.Parse(hostip); err != nil {
		return false, "", err
	}
	f.bmx.Lock()
	inblacklist := f.blacklist.Check(host)
	f.bmx.Unlock()
	if inblacklist {
		return false, "blocked by blacklist", nil
	}
	f.wmx.Lock()
	inwhitelist := f.whitelist.Check(host)
	f.wmx.Unlock()
	if inwhitelist {
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
	f.wmx.Lock()
	added := f.whitelist.Add(snet)
	f.wmx.Unlock()
	return added, nil
}

func (f *filter) DeleteWhiteList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	f.wmx.Lock()
	defer f.wmx.Unlock()
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
	f.bmx.Lock()
	added := f.blacklist.Add(snet)
	f.bmx.Unlock()
	return added, nil
}

func (f *filter) DeleteBlackList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	f.bmx.Lock()
	defer f.bmx.Unlock()
	if !f.blacklist.Exist(snet) {
		return false, nil
	}
	deleted := f.blacklist.Delete(snet)
	return deleted, nil
}

func (f *filter) GetLimits() buckets.RatesLimits {
	return f.limits
}
