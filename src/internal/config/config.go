package config

import (
	"github.com/kelseyhightower/envconfig"

	"gihub.com/bongerka/sberPDRIS/internal/config/server"
)

type Config struct {
	Server *server.Config `envconfig:"SERVER"`
}

func NewConfig[K any](prefix ...string) (*K, error) {
	var v K
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}

	if err := envconfig.Process(p, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
