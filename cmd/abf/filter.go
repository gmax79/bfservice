package main

import (
	"fmt"

	"github.com/gmax79/bfservice/internal/netsupport"
)

// filter - main objects to filtering login attempts
type filter struct {
	whitelist *netsupport.SubnetsList
	blacklist *netsupport.SubnetsList
}

// createFilter - create instance of filter
func createFilter() *filter {
	f := filter{}
	f.whitelist = netsupport.CreateSubnetsList()
	f.blacklist = netsupport.CreateSubnetsList()
	return &f
}

func (f *filter) CheckLogin(login, password, hostip string) (bool, error) {

	var host netsupport.IPAddr
	if err := host.Parse(hostip); err != nil {
		return false, err
	}
	if f.blacklist.Check(host) {
		return false, fmt.Errorf("blocked by blacklist")
	}
	if f.whitelist.Check(host) {
		return true, fmt.Errorf("passed by whitelist")
	}

	return false, nil
}

func (f *filter) ResetLogin(login, hostip string) (bool, error) {
	return false, nil
}

//todo bool remove ?
func (f *filter) AddWhiteList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	f.whitelist.Add(snet)
	return true, nil
}

func (f *filter) DeleteWhiteList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	if !f.whitelist.Exist(snet) {
		return false, nil
	}
	f.whitelist.Delete(snet)
	return true, nil
}

func (f *filter) AddBlackList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	f.blacklist.Add(snet)
	return true, nil
}

func (f *filter) DeleteBlackList(subnetip string) (bool, error) {
	var snet netsupport.Subnet
	if err := snet.Parse(subnetip); err != nil {
		return false, err
	}
	if !f.blacklist.Exist(snet) {
		return false, nil
	}
	f.blacklist.Delete(snet)
	return true, nil
}
