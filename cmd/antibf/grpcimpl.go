package main

import (
	"context"
	"net"

	grpcapi "github.com/gmax79/antibf/api/grpc"
	"github.com/gmax79/antibf/internal/buckets"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// AntibfGrpcImpl - grpc implementaion struct for service
type AntibfGrpcImpl struct {
	server    *grpc.Server
	lasterror error
	logger    *zap.Logger
	filter    buckets.Filter
}

// createGRPC - service grpc interface
func openGRPCConnect(filter buckets.Filter, host string, zaplog *zap.Logger) (*AntibfGrpcImpl, error) {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	g := &AntibfGrpcImpl{}
	g.server = grpc.NewServer()
	g.logger = zaplog
	g.filter = filter
	grpcapi.RegisterAntiBruteforceServer(g.server, g)

	go func() {
		g.lasterror = g.server.Serve(listen)
	}()

	return g, nil
}

// Stop - gracefully stopping grpc server
func (ab *AntibfGrpcImpl) Stop(ctx context.Context) {
	//todo
}

// CheckLogin - check login for bruteforce state. return true if can login or false for not
func (ab *AntibfGrpcImpl) CheckLogin(ctx context.Context,
	in *grpcapi.CheckLoginRequest) (*grpcapi.CheckLoginResponse, error) {
	var out grpcapi.CheckLoginResponse
	err := ab.filter.CheckLogin(in.Login, in.Password, in.Ip)
	out.Checked = true
	return &out, err
}

// ResetLogin - remove login from internal base (reset bruteforce rate)
func (ab *AntibfGrpcImpl) ResetLogin(ctx context.Context, in *grpcapi.ResetLoginRequest) (*grpcapi.ResetLoginResponse, error) {
	var out grpcapi.ResetLoginResponse
	err := ab.filter.ResetLogin(in.Login, in.Ip)
	out.Reseted = false
	return &out, err
}

// AddWhiteList - add ip into whitelist
func (ab *AntibfGrpcImpl) AddWhiteList(ctx context.Context, in *grpcapi.AddWhiteListRequest) (*grpcapi.AddWhiteListResponse, error) {
	var out grpcapi.AddWhiteListResponse
	err := ab.filter.AddWhiteList(in.Ipmask)
	out.Added = false
	return &out, err
}

// DeleteWhiteList - delete ip from whitelist
func (ab *AntibfGrpcImpl) DeleteWhiteList(ctx context.Context, in *grpcapi.DeleteWhiteListRequest) (*grpcapi.DeleteWhiteListResponse, error) {
	var out grpcapi.DeleteWhiteListResponse
	err := ab.filter.DeleteWhiteList(in.Ipmask)
	out.Deleted = false
	return &out, err
}

// AddBlackList - add ip into blacklist
func (ab *AntibfGrpcImpl) AddBlackList(ctx context.Context, in *grpcapi.AddBlackListRequest) (*grpcapi.AddBlackListResponse, error) {
	var out grpcapi.AddBlackListResponse
	err := ab.filter.AddBlackList(in.Ipmask)
	out.Added = false
	return &out, err
}

// DeleteBlackList - delete ip from blacklist
func (ab *AntibfGrpcImpl) DeleteBlackList(ctx context.Context, in *grpcapi.DeleteBlackListRequest) (*grpcapi.DeleteBlackListResponse, error) {
	var out grpcapi.DeleteBlackListResponse
	err := ab.filter.AddBlackList(in.Ipmask)
	out.Deleted = false
	return &out, err
}
