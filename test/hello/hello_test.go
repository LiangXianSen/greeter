package hello_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	api "github.com/LiangXianSen/greeter/api"
	"github.com/LiangXianSen/greeter/hello"
)

var (
	srv    *hello.Server
	config *hello.Config

	// HTTP server endpoint
	httpServerEndpoint = "127.0.0.1:8000"
	// GRPC server endpoint
	grpcServerEndpoint = "127.0.0.1:8080"
)

func TestMain(m *testing.M) {
	config = hello.LoadConfigByDefault()
	config.Logger.Level.SetLevel(zapcore.FatalLevel) // disable logging

	var err error
	if srv, err = hello.NewServer(config); err != nil {
		log.Fatalf("get hello service instance failed: %s", err)
	}

	go func() {
		if err := srv.ServeGRPC(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := srv.ServeHTTP(); err != nil {
			log.Fatal(err)
		}
	}()

	os.Exit(m.Run())
}

func TestHelloGRPC(t *testing.T) {
	assert := assert.New(t)
	must := require.New(t)

	conn, err := grpc.Dial(grpcServerEndpoint, grpc.WithInsecure())
	must.NoError(err)
	defer conn.Close()

	client := api.NewGreeterClient(conn)
	reply, err := client.Hello(context.Background(), &api.Request{Name: "Kevin"})
	must.NoError(err)
	assert.Equal("Hello Kevin", reply.Msg)
}

func BenchmarkHelloGRPC(b *testing.B) {
	conn, err := grpc.Dial(grpcServerEndpoint, grpc.WithInsecure())
	if err != nil {
		b.Errorf("dial grpc server failed: %s", err)
		b.FailNow()
	}
	defer conn.Close()

	client := api.NewGreeterClient(conn)
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("name-%d", i)
		_, err := client.Hello(context.Background(), &api.Request{Name: name})
		if err != nil {
			b.Errorf("%s says hello to grpc server occurred error: %s", name, err)
			b.Fail()
		}
	}
}

func TestHelloHTTP(t *testing.T) {
	assert := assert.New(t)
	must := require.New(t)

	payload := []byte(`{
		"name": "Kevin"
	}`)
	resp, err := http.Post("http://"+httpServerEndpoint+"/v1/greeter/hello", "Content-Type: application/json", bytes.NewBuffer(payload))
	must.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	must.NoError(err)

	var reply api.Response
	must.NoError(json.Unmarshal(body, &reply))
	assert.Equal("Hello Kevin", reply.Msg)
}

func BenchmarkHelloHTTP(b *testing.B) {
	payload := []byte(`{
		"name": "Kevin"
	}`)

	data := bytes.NewBuffer(payload)
	for i := 0; i < b.N; i++ {
		_, err := http.Post("http://"+httpServerEndpoint+"/v1/greeter/hello", "Content-Type: application/json", data)
		if err != nil {
			b.Errorf("Says hello to http server occurred error: %s", err)
			b.Fail()
		}
	}
}
