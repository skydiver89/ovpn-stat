package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server Server `yaml:"server"`
}

type Server struct {
	Log  string `yaml:"log"`
	Port string `yaml:"port"`
}

func (c *Config) load(fname string) error {
	yamlData, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlData, c)
	return err
}
