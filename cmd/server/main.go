package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"

	"github.com/LiangXianSen/greeter/hello"
)

var (
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
	// HTTP server endpoint
	httpServerEndpoint = flag.String("http-server-endpoint", "localhost:8080", "HTTP server endpoint")
)

func main() {
	var err error

	srv := server.NewServer(*grpcServerEndpoint, *httpServerEndpoint)

	go func() {
		if err = srv.ServeGRPC(); err != nil {
			log.Fatal(err)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err = srv.ServeHTTP(ctx, opts...); err != nil {
		log.Fatal(err)
	}
}
