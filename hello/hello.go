package hello

import (
	"context"
	"net/http"
	"sync"

	"google.golang.org/grpc"

	api "github.com/LiangXianSen/greeter/api"
)

// Server is the main struct of hello service.
type Server struct {
	config  *Config
	options *options

	httpSrv *http.Server
	grpcSrv *grpc.Server

	sg   sync.WaitGroup
	exit chan struct{}

	api.UnimplementedGreeterServer
}

// Hello says hello.
func (srv *Server) Hello(ctx context.Context, r *api.Request) (*api.Response, error) {
	reply := &api.Response{Msg: "Hello " + r.Name}
	return reply, nil
}
