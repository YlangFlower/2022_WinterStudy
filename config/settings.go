package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Type     string `yaml:"type"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Protocol string `yaml:"protocol"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Db       string `yaml:"db"`
	} `yaml:"database"`
}

func GetSettings(cfg *Config) {
	f, err := os.Open("config/settings.yaml")
	if err != nil {
		fmt.Print(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Print(err)
	}
}
