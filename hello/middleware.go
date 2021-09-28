package hello

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"

	grpc_prometheus "github.com/LiangXianSen/greeter/prometheus"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

var (
	logger *zap.Logger
)

func (srv *Server) middlewareChain() grpc.ServerOption {
	logger = srv.zapLogger()

	// interceptor chain
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_recovery.UnaryServerInterceptor(srv.options.recoveryOpts...),
		grpc_zap.UnaryServerInterceptor(logger, srv.options.zapOpts...),
		grpc_prometheus.UnaryServerInterceptor(),
	))
}

func (srv *Server) zapLogger() (logger *zap.Logger) {
	conf := zap.NewProductionConfig()
	conf.Level = srv.config.Logger.Level

	var err error
	if logger, err = conf.Build(); err != nil {
		panic(err)
	}

	return
}
