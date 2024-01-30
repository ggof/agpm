package main

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version      string   `yaml:"version"`
	Toolchain    int      `yaml:"toolchain"`
	Repositories []string `yaml:"repositories"`
	Dependencies []string `yaml:"dependencies"`
}

const (
	cfgName = "module.yaml"
)

var defaultConfig = &Config{
	Toolchain:    21,
	Repositories: []string{"maven"},
}

func main() {
	loadConfig()
}

func loadConfig() {
	cfgFile, err := os.Open(cfgName)
	if err != nil {
		log.Println("Config file not found, using default config as-is")
		return
	}
	bytes, err := io.ReadAll(cfgFile)
	if err != nil {
		log.Printf("Failed to read the config file: %s", err.Error())
		log.Println("using default config as-is")
		return
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		log.Printf("Failed to parse the config file: %s", err.Error())
		log.Println("using default config as-is")
		return
	}

	if cfg.Toolchain > 0 {
		defaultConfig.Toolchain = cfg.Toolchain
	}

	defaultConfig.Dependencies = append(defaultConfig.Dependencies, cfg.Dependencies...)
	defaultConfig.Repositories = append(defaultConfig.Repositories, cfg.Repositories...)
}
