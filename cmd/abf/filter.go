package main

import (
	"sync"
	"time"

	"github.com/gmax79/bfservice/internal/netsupport"
	"github.com/gmax79/bfservice/internal/ratelimit"
	"github.com/gmax79/bfservice/internal/storage"
	"github.com/jdeal-mediamath/clockwork"
)

// filter - main objects to filtering login attempts
type filter struct {
	whitelist *netsupport.SubnetsList
	blacklist *netsupport.SubnetsList
	wmx, bmx  *sync.Mutex
	counter   *ratelimit.Counter
	limits    ratelimit.Config
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
	f.limits.LoginDuration = time.Minute
	f.limits.Password = config.PasswordRate
	f.limits.PasswordDuration = time.Minute
	f.limits.Host = config.IPRate
	f.limits.HostDuration = time.Minute

	f.counter, err = ratelimit.CreateCounter(f.limits, clockwork.NewRealClock())
	if err != nil {
		return nil, err
	}
	f.wmx = &sync.Mutex{}
	f.bmx = &sync.Mutex{}
	return &f, nil
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
	scored, reason := f.counter.CheckAndCount(login, password, hostip)
	return scored, reason, nil
}

func (f *filter) ResetLogin(login, hostip string) bool {
	return f.counter.Reset(login, hostip)
}

func (f *filter) AddWhiteList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	f.wmx.Lock()
	added, err := f.whitelist.Add(snet)
	f.wmx.Unlock()
	return added, err
}

func (f *filter) DeleteWhiteList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	f.wmx.Lock()
	deleted, err := f.whitelist.Delete(snet)
	f.wmx.Unlock()
	return deleted, err
}

func (f *filter) AddBlackList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	f.bmx.Lock()
	added, err := f.blacklist.Add(snet)
	f.bmx.Unlock()
	return added, err
}

func (f *filter) DeleteBlackList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	f.bmx.Lock()
	deleted, err := f.blacklist.Delete(snet)
	f.bmx.Unlock()
	return deleted, err
}

func (f *filter) GetLimits() ratelimit.Config {
	return f.limits
}

func (f *filter) Close() error {
	return f.stor.Close()
}
