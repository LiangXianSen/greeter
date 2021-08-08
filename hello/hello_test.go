package hello

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/LiangXianSen/greeter/api"
)

func TestHelloGRPC(t *testing.T) {
	assert := assert.New(t)
	must := require.New(t)
	req := &api.Request{
		Name: "Kevin",
	}
	resp, err := testSrv.Hello(context.Background(), req)
	must.NoError(err)
	assert.Equal("Hello Kevin", resp.Msg)
}

func BenchmarkHelloGRPC(b *testing.B) {
	req := &api.Request{
		Name: "Kevin",
	}
	for i := 0; i < b.N; i++ {
		if _, err := testSrv.Hello(context.Background(), req); err != nil {
			b.Errorf("says hello to grpc server occurred error: %s", err)
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

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/v1/greeter/hello", bytes.NewReader(payload))

	testSrv.httpSrv.Handler.ServeHTTP(w, r)

	body, err := ioutil.ReadAll(w.Body)
	must.NoError(err)

	var reply api.Response
	must.NoError(json.Unmarshal(body, &reply))
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("Hello Kevin", reply.Msg)
}

func BenchmarkHelloHTTP(b *testing.B) {
	payload := []byte(`{
		"name": "Kevin"
	}`)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/v1/greeter/hello", bytes.NewReader(payload))

	for i := 0; i < b.N; i++ {
		testSrv.httpSrv.Handler.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			b.Errorf("says hello to http server occurred error: %d", w.Code)
			b.Fail()
		}
	}
}
