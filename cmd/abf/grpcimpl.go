package main

import (
	"context"
	"net"

	grpcapi "github.com/gmax79/bfservice/api/grpc"
	"github.com/gmax79/bfservice/internal/buckets"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// AbfGrpcImpl - grpc implementaion struct for service
type AbfGrpcImpl struct {
	server    *grpc.Server
	lasterror error
	logger    *zap.Logger
	filter    buckets.Filter
}

// openGRPCServer - service grpc interface
func openGRPCServer(filter buckets.Filter, host string, zaplog *zap.Logger) (*AbfGrpcImpl, error) {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}
	g := &AbfGrpcImpl{}
	g.server = grpc.NewServer()
	g.logger = zaplog
	g.filter = filter
	grpcapi.RegisterAntiBruteforceServer(g.server, g)

	go func() {
		g.lasterror = g.server.Serve(listen)
	}()
	return g, nil
}

// HealthCheck - method to check service for alive
func (ab *AbfGrpcImpl) HealthCheck(ctx context.Context, in *grpcapi.HealthCheckRequst) (*grpcapi.HealthCheckResponse, error) {
	var out grpcapi.HealthCheckResponse
	out.Status = "ok"
	return &out, nil
}

// Stop - gracefully stopping grpc server
func (ab *AbfGrpcImpl) Stop(ctx context.Context) {
	ab.server.GracefulStop()
	//todo use ctx
}

// CheckLogin - check login for bruteforce state. return true if can login or false for not
func (ab *AbfGrpcImpl) CheckLogin(ctx context.Context,
	in *grpcapi.CheckLoginRequest) (*grpcapi.CheckLoginResponse, error) {
	var out grpcapi.CheckLoginResponse
	err := ab.filter.CheckLogin(in.Login, in.Password, in.Ip)
	out.Checked = true
	out.Reason = "Not limited"
	return &out, err
}

// ResetLogin - remove login from internal base (reset bruteforce rate)
func (ab *AbfGrpcImpl) ResetLogin(ctx context.Context, in *grpcapi.ResetLoginRequest) (*grpcapi.ResetLoginResponse, error) {
	var out grpcapi.ResetLoginResponse
	err := ab.filter.ResetLogin(in.Login, in.Ip)
	out.Reseted = false
	return &out, err
}

// AddWhiteList - add ip into whitelist
func (ab *AbfGrpcImpl) AddWhiteList(ctx context.Context, in *grpcapi.AddWhiteListRequest) (*grpcapi.AddWhiteListResponse, error) {
	var out grpcapi.AddWhiteListResponse
	err := ab.filter.AddWhiteList(in.Ipmask)
	out.Added = false
	return &out, err
}

// DeleteWhiteList - delete ip from whitelist
func (ab *AbfGrpcImpl) DeleteWhiteList(ctx context.Context, in *grpcapi.DeleteWhiteListRequest) (*grpcapi.DeleteWhiteListResponse, error) {
	var out grpcapi.DeleteWhiteListResponse
	err := ab.filter.DeleteWhiteList(in.Ipmask)
	out.Deleted = false
	return &out, err
}

// AddBlackList - add ip into blacklist
func (ab *AbfGrpcImpl) AddBlackList(ctx context.Context, in *grpcapi.AddBlackListRequest) (*grpcapi.AddBlackListResponse, error) {
	var out grpcapi.AddBlackListResponse
	err := ab.filter.AddBlackList(in.Ipmask)
	out.Added = false
	return &out, err
}

// DeleteBlackList - delete ip from blacklist
func (ab *AbfGrpcImpl) DeleteBlackList(ctx context.Context, in *grpcapi.DeleteBlackListRequest) (*grpcapi.DeleteBlackListResponse, error) {
	var out grpcapi.DeleteBlackListResponse
	err := ab.filter.AddBlackList(in.Ipmask)
	out.Deleted = false
	return &out, err
}
