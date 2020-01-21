package main

import (
	"context"
	"net"

	grpcapi "github.com/gmax79/antibf/api/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// AntibfGrpcImpl - grpc implementaion struct for service
type AntibfGrpcImpl struct {
	server    *grpc.Server
	lasterror error
	logger    *zap.Logger
}

// createGRPC - service grpc interface
func openGRPCConnect(host string, zaplog *zap.Logger) (*AntibfGrpcImpl, error) {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	g := &AntibfGrpcImpl{}
	g.server = grpc.NewServer()
	g.logger = zaplog
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
	var res grpcapi.CheckLoginResponse
	return &res, nil
}

// ResetLogin - remove login from internal base (reset bruteforce rate)
func (ab *AntibfGrpcImpl) ResetLogin(ctx context.Context, in *grpcapi.ResetLoginRequest) (*grpcapi.ResetLoginResponse, error) {
	var res grpcapi.ResetLoginResponse
	return &res, nil
}

// AddWhiteList - add ip into whitelist
func (ab *AntibfGrpcImpl) AddWhiteList(ctx context.Context, in *grpcapi.AddWhiteListRequest) (*grpcapi.AddWhiteListResponse, error) {
	var res grpcapi.AddWhiteListResponse
	return &res, nil
}

// DeleteWhiteList - delete ip from whitelist
func (ab *AntibfGrpcImpl) DeleteWhiteList(ctx context.Context, in *grpcapi.DeleteWhiteListRequest) (*grpcapi.DeleteWhiteListResponse, error) {
	var res grpcapi.DeleteWhiteListResponse
	return &res, nil
}

// AddBlackList - add ip into blacklist
func (ab *AntibfGrpcImpl) AddBlackList(ctx context.Context, in *grpcapi.AddBlackListRequest) (*grpcapi.AddBlackListResponse, error) {
	var res grpcapi.AddBlackListResponse
	return &res, nil
}

// DeleteBlackList - delete ip from blacklist
func (ab *AntibfGrpcImpl) DeleteBlackList(ctx context.Context, in *grpcapi.DeleteBlackListRequest) (*grpcapi.DeleteBlackListResponse, error) {
	var res grpcapi.DeleteBlackListResponse
	return &res, nil
}
