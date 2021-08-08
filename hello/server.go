package hello

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	api "github.com/LiangXianSen/greeter/api"
)

// NewServer returns the hello service instance.
func NewServer(conf *Config, opt ...Option) (srv *Server, err error) {
	opts := new(options)
	*opts = *defaultOptions
	for _, o := range opt {
		o(opts)
	}

	srv = &Server{
		config:  conf,
		options: opts,
		exit:    make(chan struct{}),
	}

	// grpc
	md := srv.middlewareChain()
	grpcServer := grpc.NewServer(md)
	srv.grpcSrv = grpcServer
	api.RegisterGreeterServer(grpcServer, srv)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// http gateway
	mux := runtime.NewServeMux()
	if err = api.RegisterGreeterHandlerServer(ctx, mux, srv); err != nil {
		return
	}
	listener := &http.Server{
		Handler: mux,
		Addr:    srv.config.HTTPEndpoint.String(),
	}
	srv.httpSrv = listener

	return srv, nil
}

// ServeHTTP serves HTTP endpoint.
func (srv *Server) ServeHTTP() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	go func() {
		<-srv.exit
		defer srv.sg.Done()
		srv.httpSrv.Shutdown(ctx)
	}()

	srv.sg.Add(1)
	return srv.httpSrv.ListenAndServe()
}

// ServeGRPC serves GRPC endpoint.
func (srv *Server) ServeGRPC() (err error) {
	var lis net.Listener
	if lis, err = net.Listen("tcp", srv.config.GRPCEndpoint.String()); err != nil {
		return
	}

	go func() {
		<-srv.exit
		defer srv.sg.Done()
		srv.grpcSrv.GracefulStop()
	}()

	srv.sg.Add(1)
	return srv.grpcSrv.Serve(lis)
}

// Shutdown shuts the server down.
func (srv *Server) Shutdown() {
	close(srv.exit)
	srv.sg.Wait()
}

// SetLoggerLevel sets logger level.
func (srv *Server) SetLoggerLevel(level zapcore.Level) {
	srv.config.Logger.Level.SetLevel(level)
}
