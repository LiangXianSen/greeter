package prometheus

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		serviceName, methodName := splitMethodName(info.FullMethod)
		timer := prometheus.NewTimer(prometheus.ObserverFunc(requestDuration.WithLabelValues("unary", serviceName, methodName).Observe))
		defer timer.ObserveDuration()

		resp, err := handler(ctx, req)

		st, _ := status.FromError(err)
		requestCounter.WithLabelValues("unary", serviceName, methodName, st.Code().String()).Inc()
		return resp, err
	}
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		err = handler(srv, stream)
		return err
	}
}
