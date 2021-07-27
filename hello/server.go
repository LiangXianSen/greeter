package server

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	api "github.com/LiangXianSen/greeter/api"
)

type HelloService struct {
	httpServerEndpoint string
	grpcServerEndpoint string
	api.UnimplementedGreeterServer
}

func NewServer(grpcServerEndpoint, httpServerEndpoint string) *HelloService {
	helloSrv := &HelloService{
		httpServerEndpoint: httpServerEndpoint,
		grpcServerEndpoint: grpcServerEndpoint,
	}
	return helloSrv
}

func (srv *HelloService) Hello(ctx context.Context, r *api.Request) (*api.Response, error) {
	reply := &api.Response{Msg: "Hello: " + r.Name}
	return reply, nil
}

func (srv *HelloService) ServeHTTP(ctx context.Context, opts ...grpc.DialOption) (err error) {
	mux := runtime.NewServeMux()
	if err = api.RegisterGreeterHandlerFromEndpoint(ctx, mux, srv.grpcServerEndpoint, opts); err != nil {
		return
	}
	return http.ListenAndServe(srv.httpServerEndpoint, mux)
}

func (srv *HelloService) ServeGRPC() (err error) {
	grpcServer := grpc.NewServer()
	api.RegisterGreeterServer(grpcServer, srv)

	var lis net.Listener
	if lis, err = net.Listen("tcp", srv.grpcServerEndpoint); err != nil {
		return
	}

	return grpcServer.Serve(lis)
}
