package tests

import (
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Scheme       string `required:"true"`
	Host         string `required:"true"`
	Port         *int
	EnvNamespace *string
}

func LoadConfig(envPrefix string) (*Config, error) {
	config := Config{}
	if err := envconfig.Process(envPrefix, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) BuildURL(endpoint string) (*string, error) {
	if len(endpoint) == 0 {
		return nil, errors.New("empty endpoint is not valid")
	}
	if string(endpoint[0]) != "/" {
		return nil, errors.New("endpoint must start with /")
	}

	url := fmt.Sprintf("%s://%s", c.Scheme, c.Host)
	if c.Port != nil {
		url = fmt.Sprintf("%s:%d", url, *c.Port)
	}
	if c.EnvNamespace != nil {
		url = fmt.Sprintf("%s/%s", url, *c.EnvNamespace)
	}
	url = fmt.Sprintf("%s%s", url, endpoint)
	return &url, nil
}
