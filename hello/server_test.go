package hello

import (
	"log"
	"os"
	"testing"

	"go.uber.org/zap/zapcore"
)

var (
	testSrv *Server
	config  *Config
)

func TestMain(m *testing.M) {
	config = LoadConfigByDefault()
	config.Logger.Level.SetLevel(zapcore.FatalLevel) // disable logging

	var err error
	if testSrv, err = NewServer(config); err != nil {
		log.Fatalf("get hello service instance failed: %s", err)
	}

	go func() {
		if err := testSrv.ServeGRPC(); err != nil {
			log.Fatal(err)
		}
	}()
	testSrv.IsGRPCReady()

	os.Exit(m.Run())

	// teardown
	testSrv.Shutdown()
}
