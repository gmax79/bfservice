package main

import (
	"context"
	"net"

	grpcapi "github.com/gmax79/bfservice/api/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// AbfGrpcImpl - grpc implementation struct for service
type AbfGrpcImpl struct {
	server    *grpc.Server
	lasterror error
	logger    *zap.Logger
	hfilter   *filter
}

// openGRPCServer - service grpc interface
func openGRPCServer(config RatesAndHostConfig, zaplog *zap.Logger) (*AbfGrpcImpl, error) {
	listen, err := net.Listen("tcp", config.Host)
	if err != nil {
		return nil, err
	}
	var g AbfGrpcImpl
	g.hfilter, err = createFilter(config)
	if err != nil {
		return nil, err
	}
	g.server = grpc.NewServer()
	g.logger = zaplog
	grpcapi.RegisterAntiBruteforceServer(g.server, &g)

	go func() {
		g.lasterror = g.server.Serve(listen)
	}()
	return &g, nil
}

// HealthCheck - method to check service for alive
func (ab *AbfGrpcImpl) HealthCheck(ctx context.Context, in *grpcapi.HealthCheckRequst) (*grpcapi.HealthCheckResponse, error) {
	var out grpcapi.HealthCheckResponse
	out.Status = "ok"
	return &out, nil
}

// Stop - gracefully stopping grpc server
func (ab *AbfGrpcImpl) Stop() {
	ab.server.GracefulStop()
}

// CheckLogin - check login for bruteforce state. return true if can login or false for not
func (ab *AbfGrpcImpl) CheckLogin(ctx context.Context, in *grpcapi.CheckLoginRequest) (*grpcapi.CheckLoginResponse, error) {
	var out grpcapi.CheckLoginResponse
	var err error
	out.Checked, out.Reason, err = ab.hfilter.CheckLogin(in.Login, in.Password, in.Ip)
	// Skip log per call. Main func with expensive rate
	return &out, err
}

// ResetLogin - remove login from internal base (reset bruteforce rate)
func (ab *AbfGrpcImpl) ResetLogin(ctx context.Context, in *grpcapi.ResetLoginRequest) (*grpcapi.ResetLoginResponse, error) {
	var out grpcapi.ResetLoginResponse
	out.Reseted = ab.hfilter.ResetLogin(in.Login, in.Ip)
	if out.Reseted {
		out.Reason = "exists"
	} else {
		out.Reason = "not found"
	}
	ab.logger.Info("Reset login/ip", zap.String("login", in.Login), zap.String("host", in.Ip), zap.Bool("was exist", out.Reseted))
	return &out, nil
}

// AddWhiteList - add ip into whitelist
func (ab *AbfGrpcImpl) AddWhiteList(ctx context.Context, in *grpcapi.AddWhiteListRequest) (*grpcapi.AddWhiteListResponse, error) {
	var out grpcapi.AddWhiteListResponse
	var err error
	out.Added, err = ab.hfilter.AddWhiteList(in.Ipmask)
	if err != nil {
		ab.logger.Error("Add in whitelist", zap.String("mask", in.Ipmask), zap.Error(err))
		out.Reason = err.Error()
	} else {
		ab.logger.Info("Add in whitelist", zap.String("mask", in.Ipmask), zap.Bool("already exist", !out.Added))
		if !out.Added {
			out.Reason = "already exist"
		} else {
			out.Reason = "new element"
		}
	}
	return &out, err
}

// DeleteWhiteList - delete ip from whitelist
func (ab *AbfGrpcImpl) DeleteWhiteList(ctx context.Context, in *grpcapi.DeleteWhiteListRequest) (*grpcapi.DeleteWhiteListResponse, error) {
	var out grpcapi.DeleteWhiteListResponse
	var err error
	out.Deleted, err = ab.hfilter.DeleteWhiteList(in.Ipmask)
	if err != nil {
		ab.logger.Error("Delete from whitelist", zap.String("mask", in.Ipmask), zap.Error(err))
		out.Reason = err.Error()
	} else {
		ab.logger.Info("Delete from whitelist", zap.String("mask", in.Ipmask), zap.Bool("was exist", out.Deleted))
		if !out.Deleted {
			out.Reason = "not exist"
		} else {
			out.Reason = "deleted"
		}
	}
	return &out, err
}

// AddBlackList - add ip into blacklist
func (ab *AbfGrpcImpl) AddBlackList(ctx context.Context, in *grpcapi.AddBlackListRequest) (*grpcapi.AddBlackListResponse, error) {
	ab.logger.Info("Add in blacklist", zap.String("mask", in.Ipmask))
	var out grpcapi.AddBlackListResponse
	var err error
	out.Added, err = ab.hfilter.AddBlackList(in.Ipmask)
	if err != nil {
		ab.logger.Error("Add in blacklist", zap.String("mask", in.Ipmask), zap.Error(err))
		out.Reason = err.Error()
	} else {
		ab.logger.Info("Add in blacklist", zap.String("mask", in.Ipmask), zap.Bool("already exist", !out.Added))
		if !out.Added {
			out.Reason = "already exist"
		} else {
			out.Reason = "new element"
		}
	}
	return &out, err
}

// DeleteBlackList - delete ip from blacklist
func (ab *AbfGrpcImpl) DeleteBlackList(ctx context.Context, in *grpcapi.DeleteBlackListRequest) (*grpcapi.DeleteBlackListResponse, error) {
	ab.logger.Info("Delete from blacklist", zap.String("mask", in.Ipmask))
	var out grpcapi.DeleteBlackListResponse
	var err error
	out.Deleted, err = ab.hfilter.DeleteBlackList(in.Ipmask)
	if err != nil {
		ab.logger.Error("Delete from blacklist", zap.String("mask", in.Ipmask), zap.Error(err))
		out.Reason = err.Error()
	} else {
		ab.logger.Info("Delete from blacklist", zap.String("mask", in.Ipmask), zap.Bool("was exist", out.Deleted))
		if !out.Deleted {
			out.Reason = "not exist"
		} else {
			out.Reason = "deleted"
		}
	}
	return &out, err
}

// GetRates - get current config settings of service
func (ab *AbfGrpcImpl) GetRates(ctx context.Context, in *grpcapi.GetRatesRequest) (*grpcapi.GetRatesResponse, error) {
	var out grpcapi.GetRatesResponse
	limits := ab.hfilter.GetLimits()
	out.LoginRate = int32(limits.Login)
	out.LoginInterval = limits.LoginDuration.Milliseconds()
	out.PasswordRate = int32(limits.Password)
	out.PasswordInterval = limits.PasswordDuration.Milliseconds()
	out.HostRate = int32(limits.Host)
	out.HostInterval = limits.HostDuration.Milliseconds()
	return &out, nil
}
