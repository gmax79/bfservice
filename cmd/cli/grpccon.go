package main

import (
	"context"
	"errors"

	grpcapi "github.com/gmax79/antibf/api/grpc"

	"google.golang.org/grpc"
)

// grpcConn - connection to service
type grpcConn struct {
	cancel func()
	client grpcapi.AntiBruteforceClient
}

func grpcConnect(host string) (*grpcConn, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	// var err error
	// cli, err :=     grpccon.CreateClient(ctx, host)
	//	assert("", err)

	clientCon, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := &grpcConn{}
	// var cancelfunc func()
	// _, cancelfunc = context.WithCancel(context.Background())
	c.cancel = func() {
		// cancelfunc()
		clientCon.Close()
	}
	c.client = grpcapi.NewAntiBruteforceClient(clientCon)
	return c, nil
}

func (c *grpcConn) HealthCheck(ctx context.Context) error {
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
