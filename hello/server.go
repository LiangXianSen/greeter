package hello

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	srv.grpcSrv = srv.newGRPCService()

	if srv.httpSrv, err = srv.newHTTPGateway(); err != nil {
		return
	}
	return
}

func (srv *Server) newGRPCService() *grpc.Server {
	md := srv.middlewareChain()
	grpcServer := grpc.NewServer(md)
	api.RegisterGreeterServer(grpcServer, srv)

	return grpcServer
}

func (srv *Server) newHTTPGateway() (gateway *http.Server, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(
		srv.config.GRPCEndpoint.String(),
		grpc.WithInsecure(),
	); err != nil {
		return
	}

	mux := runtime.NewServeMux()
	if err = api.RegisterGreeterHandler(context.Background(), mux, conn); err != nil {
		return
	}
	gateway = &http.Server{
		Handler: mux,
		Addr:    srv.config.HTTPEndpoint.String(),
	}

	// metrics
	mux.HandlePath("GET", "/metrics", metricsHandler)
	return
}

func metricsHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	promhttp.Handler().ServeHTTP(w, r)
}

// ServeHTTP serves HTTP endpoint.
func (srv *Server) ServeHTTP() (err error) {
	ctx := context.Background()

	go func() {
		<-srv.exit
		defer srv.sg.Done()
		srv.httpSrv.Shutdown(ctx)
	}()

	srv.sg.Add(1)
	log.Printf("HTTP service listen on :%d", srv.config.HTTPEndpoint.Port)
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
	log.Printf("GRPC service listen on :%d", srv.config.GRPCEndpoint.Port)
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

// IsGRPCReady checks the GRPC socket status.
func (srv *Server) IsGRPCReady() bool {
	var i, limit = 0, 10
	for i = 0; i < limit; i++ {
		conn, err := net.Dial("tcp", srv.config.GRPCEndpoint.String())
		if err != nil {
			time.Sleep(time.Second * 1)
			continue
		}
		conn.Close()
		break
	}
	return true && i < limit
}

// IsHTTPReady checks HTTP socket status.
func (srv *Server) IsHTTPReady() bool {
	var i, limit = 0, 10
	for i = 0; i < limit; i++ {
		conn, err := net.Dial("tcp", srv.config.HTTPEndpoint.String())
		if err != nil {
			time.Sleep(time.Second * 1)
			continue
		}
		conn.Close()
		break
	}
	return true && i < limit
}
