package hello

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)
	must := require.New(t)
	data := bytes.NewReader([]byte(`{
		"http_endpoint": {
			"address": "127.0.0.1",
			"port": 1234
		},
		"grpc_endpoint": {
			"address": "0.0.0.0",
			"port": 1111
		}
	}`))

	conf, err := LoadConfig(data, "json")
	must.NoError(err)
	assert.Equal("127.0.0.1", conf.HTTPEndpoint.Address)
	assert.Equal(uint16(1234), conf.HTTPEndpoint.Port)
	assert.Equal("0.0.0.0", conf.GRPCEndpoint.Address)
	assert.Equal(uint16(1111), conf.GRPCEndpoint.Port)
}

func TestLoadConfigFromFile(t *testing.T) {
	assert := assert.New(t)
	must := require.New(t)

	conf, err := LoadConfigFromFile("../config/config.json")
	must.NoError(err)
	assert.Equal(defaultHTTPPort, conf.HTTPEndpoint.Port)
	assert.Equal(defaultGRPCPort, conf.GRPCEndpoint.Port)
}
