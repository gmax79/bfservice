package grpccon

import (
	"context"
	"errors"
	"time"

	grpcapi "github.com/gmax79/bfservice/api/grpc"

	"google.golang.org/grpc"
)

// Response - answer from service, if error not occur
type Response struct {
	Status bool
	Reason string
}

// Rates - struct with current settings and state
type Rates struct {
	LoginRate        int
	LoginInterval    time.Duration
	PasswordRate     int
	PasswordInterval time.Duration
	HostRate         int
	HostInterval     time.Duration
}

// Client - client to connect to service
type Client struct {
	cancel func()
	client grpcapi.AntiBruteforceClient
}

// Connect - create connection to abf service
func Connect(host string) (*Client, error) {
	clientCon, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := &Client{}
	c.cancel = func() {
		clientCon.Close()
	}
	c.client = grpcapi.NewAntiBruteforceClient(clientCon)
	return c, nil
}

// Close - disconnect from service
func (c *Client) Close() {
	c.cancel()
}

// HealthCheck - check service connect and alive
func (c *Client) HealthCheck(ctx context.Context) error {
	var req grpcapi.HealthCheckRequst
	req.Apiversion = "1"
	resp, err := c.client.HealthCheck(ctx, &req)
	if err != nil {
		return err
	}
	if resp.Status != "ok" {
		return errors.New("invalid answer from remote service")
	}
	return nil
}

// CheckLogin - check login, pass and ip for bruteforce status
func (c *Client) CheckLogin(ctx context.Context, login, password, ip string) (*Response, error) {
	var req grpcapi.CheckLoginRequest
	req.Login = login
	req.Password = password
	req.Ip = ip
	resp, err := c.client.CheckLogin(ctx, &req)
	if err != nil {
		return nil, err
	}
	var r Response
	r.Status = resp.Checked
	r.Reason = resp.Reason
	return &r, nil
}

// ResetLogin - check login, pass and ip for bruteforce status
func (c *Client) ResetLogin(ctx context.Context, login, ip string) (*Response, error) {
	var req grpcapi.ResetLoginRequest
	req.Login = login
	req.Ip = ip
	resp, err := c.client.ResetLogin(ctx, &req)
	if err != nil {
		return nil, err
	}
	var r Response
	r.Status = resp.Reseted
	r.Reason = resp.Reason
	return &r, nil
}

// AddWhiteList - add ip or subnet into whitelist
func (c *Client) AddWhiteList(ctx context.Context, ipmask string) (*Response, error) {
	var req grpcapi.AddWhiteListRequest
	req.Ipmask = ipmask
	resp, err := c.client.AddWhiteList(ctx, &req)
	if err != nil {
		return nil, err
	}
	var r Response
	r.Status = resp.Added
	r.Reason = resp.Reason
	return &r, nil
}

// DeleteWhiteList - delete ip or subnet from whitelist
func (c *Client) DeleteWhiteList(ctx context.Context, ipmask string) (*Response, error) {
	var req grpcapi.DeleteWhiteListRequest
	req.Ipmask = ipmask
	resp, err := c.client.DeleteWhiteList(ctx, &req)
	if err != nil {
		return nil, err
	}
	var r Response
	r.Status = resp.Deleted
	r.Reason = resp.Reason
	return &r, nil
}

// AddBlackList - add ip or subnet into blacklist
func (c *Client) AddBlackList(ctx context.Context, ipmask string) (*Response, error) {
	var req grpcapi.AddBlackListRequest
	req.Ipmask = ipmask
	resp, err := c.client.AddBlackList(ctx, &req)
	if err != nil {
		return nil, err
	}
	var r Response
	r.Status = resp.Added
	r.Reason = resp.Reason
	return &r, nil
}

// DeleteBlackList - delete ip or subnet from blacklist
func (c *Client) DeleteBlackList(ctx context.Context, ipmask string) (*Response, error) {
	var req grpcapi.DeleteBlackListRequest
	req.Ipmask = ipmask
	resp, err := c.client.DeleteBlackList(ctx, &req)
	if err != nil {
		return nil, err
	}
	var r Response
	r.Status = resp.Deleted
	r.Reason = resp.Reason
	return &r, nil
}

// GetRates - get current remote server rates
func (c *Client) GetRates(ctx context.Context) (*Rates, error) {
	var req grpcapi.GetRatesRequest
	resp, err := c.client.GetRates(ctx, &req)
	if err != nil {
		return nil, err
	}
	var r Rates
	r.LoginRate = int(resp.LoginRate)
	r.LoginInterval = time.Duration(resp.LoginInterval)
	r.PasswordRate = int(resp.PasswordRate)
	r.PasswordInterval = time.Duration(resp.PasswordInterval)
	r.HostRate = int(resp.HostRate)
	r.HostInterval = time.Duration(resp.HostInterval)
	return &r, nil
}
