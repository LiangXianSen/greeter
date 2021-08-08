package hello

import (
	"fmt"
	"io"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config indicates which values the server need.
type Config struct {
	Logger       *Logger   `json:"logger" mapstructure:"logger"`
	HTTPEndpoint *Endpoint `json:"http_endpoint" mapstructure:"http_endpoint"`
	GRPCEndpoint *Endpoint `json:"grpc_endpoint" mapstructure:"grpc_endpoint"`
}

// Logger indicates logging configuration.
type Logger struct {
	Level zap.AtomicLevel `json:"level" mapstructure:"level"`
}

// Endpoint indicates network listener configuration.
type Endpoint struct {
	Address string `json:"address" mapstructure:"address"`
	Port    uint16 `json:"port" mapstructure:"port"`
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Address, e.Port)
}

const (
	defaultHTTPPort uint16 = 8000
	defaultGRPCPort uint16 = 8080
)

// LoadConfigByDefault loads config by default
func LoadConfigByDefault() *Config {
	return &Config{
		Logger: &Logger{
			Level: zap.NewAtomicLevelAt(zapcore.InfoLevel),
		},
		HTTPEndpoint: &Endpoint{
			Address: "0.0.0.0",
			Port:    defaultHTTPPort,
		},
		GRPCEndpoint: &Endpoint{
			Address: "0.0.0.0",
			Port:    defaultGRPCPort,
		},
	}
}

// LoadConfigFromFile loads config from the named file path.
// Finds config file from 'config/' directory If path is empty.
func LoadConfigFromFile(path string) (conf *Config, err error) {
	vc := viper.New()
	vc.AddConfigPath(".")
	vc.AddConfigPath("config")
	vc.SetConfigName("config")

	vc.SetConfigFile(path)

	vc.SetDefault("http_endpoint.port", defaultHTTPPort)
	vc.SetDefault("grpc_endpoint.port", defaultGRPCPort)
	vc.SetDefault("logger.level", zap.NewAtomicLevelAt(zapcore.InfoLevel))

	vc.SetEnvPrefix("greeter")
	vc.AutomaticEnv()

	if err = vc.ReadInConfig(); err != nil {
		return
	}

	if err = vc.Unmarshal(&conf); err != nil {
		return
	}

	return
}

// LoadConfig loads config from io.reader. extension includes: yaml, json, toml,
// ini, hcl
func LoadConfig(fd io.Reader, extension string) (conf *Config, err error) {
	vc := viper.New()
	vc.SetConfigType(extension)

	vc.SetDefault("http_endpoint.port", defaultHTTPPort)
	vc.SetDefault("grpc_endpoint.port", defaultGRPCPort)
	vc.SetDefault("logger.level", zap.NewAtomicLevelAt(zapcore.InfoLevel))

	if err = vc.ReadConfig(fd); err != nil {
		return
	}

	if err = vc.Unmarshal(&conf); err != nil {
		return
	}

	return
}
