package server

import (
	"strings"
)

type Config struct {
	HTTPAddr string `envconfig:"HTTP_ADDR" required:"true"`
	Env      string `envconfig:"ENV" required:"true"`
}

func (c *Config) IsDev() bool {
	if strings.EqualFold(c.Env, "dev") {
		return true
	}

	return false
}
