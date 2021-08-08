package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap/zapcore"

	"github.com/LiangXianSen/greeter/hello"
)

var (
	cfgFile string
	verbose bool
)

func run() {
	var err error
	var conf *hello.Config
	if conf, err = hello.LoadConfigFromFile(cfgFile); err != nil {
		log.Fatalf("load config failed: %s", err)
	}

	if verbose {
		conf.Logger.Level.SetLevel(zapcore.DebugLevel)
	}

	var srv *hello.Server
	if srv, err = hello.NewServer(conf); err != nil {
		log.Fatalf("get hello service instance failed: %s", err)
	}

	// Serve GRPC
	go func() {
		if err = srv.ServeGRPC(); err != nil {
			log.Fatal(err)
		}
	}()

	// Serve HTTP
	go func() {
		if err = srv.ServeHTTP(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	srv.Shutdown()
	log.Println("shutdown!")
}
