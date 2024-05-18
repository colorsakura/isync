package config

import (
	"github.com/BurntSushi/toml"
)

type Server struct {
	Name     string `toml:"name"`
	Address  string `toml:"address"`
	Account  string `toml:"account"`
	Password string `toml:"password"`
}

type Dir struct {
	Directory string `toml:"directory"`
	Target    string `toml:"target"`
}

type Config struct {
	Server
	Dir
}

func UmarshalConfig(buf []byte) (*Config, error) {
	cfg := new(Config)

	if err := toml.Unmarshal(buf, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) String() string {
	return "&Config{}"
}
