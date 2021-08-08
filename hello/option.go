package hello

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ctxzap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

// Option gives optional arguments.
type Option func(*options)

type options struct {
	recoveryOpts   []grpc_recovery.Option
	zapOpts        []grpc_zap.Option
}

var defaultOptions = &options{
	recoveryOpts: []grpc_recovery.Option{grpc_recovery.WithRecoveryHandler(defaultRecoveryHandlerFunc)},
	zapOpts: []grpc_zap.Option{
		grpc_zap.WithDurationField(grpc_zap.DurationToDurationField),
		grpc_zap.WithMessageProducer(defaultMessageProducer),
	},
}

// OptRecoveryHandler is the option which give a func to handle recovery.
func OptRecoveryHandler(fn func(r interface{}) (err error)) Option {
	return func(o *options) {
		o.recoveryOpts = append(o.recoveryOpts, grpc_recovery.WithRecoveryHandler(fn))
	}
}

// OptMessageProducer is the option which give a logging message producer.
func OptMessageProducer(fn func(ctx context.Context, msg string, level zapcore.Level, code codes.Code, err error, duration zapcore.Field)) Option {
	return func(o *options) {
		o.zapOpts = append(o.zapOpts, grpc_zap.WithMessageProducer(fn))
	}
}

func defaultRecoveryHandlerFunc(r interface{}) (err error) {
	log.Printf("recover from a critical error: %v", r)
	debug.PrintStack()
	return status.Errorf(codes.Internal, "%v", r)
}

func defaultMessageProducer(ctx context.Context, msg string, level zapcore.Level, code codes.Code, err error, duration zapcore.Field) {
	ctxzap.Extract(ctx).Check(level, msg).Write(
		zap.Error(err),
		zap.String("grpc.code", code.String()),
		duration,
	)
}

func durationToTimeMillisField(duration time.Duration) zapcore.Field {
	return zap.Float32("grpc.duration_ms", durationToMilliseconds(duration))
}

func durationToMilliseconds(duration time.Duration) float32 {
	return float32(duration.Nanoseconds()/1000) / 1000
}
