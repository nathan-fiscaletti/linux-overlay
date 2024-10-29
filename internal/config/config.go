package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port    int `yaml:"port"`
	Monitor struct {
		Keys         []uint16 `yaml:"keys"`
		MouseButtons []uint16 `yaml:"mouse_buttons"`
	}
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
