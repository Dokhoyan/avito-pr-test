package config

import (
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type httpConfig struct {
	host string
	port string
}

func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if host == "" {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(httpPortEnvName)
	if port == "" {
		return nil, errors.New("http port not found")
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
