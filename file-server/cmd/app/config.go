package main

import (
	"encoding/json"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	NasRootPath string `json:"nas_root_path" yaml:"nas_root_path"`
	ScanOption  struct {
		Path       []string `json:"path" yaml:"path"`
		RegexPath  []string `json:"regex_path" yaml:"regex_path"`
		Extensions []string `json:"extensions" yaml:"extensions"`
	} `json:"scan_option" yaml:"scan_option"`
}

func (c *Config) Load(path string) {
	// load config from file
	if strings.HasSuffix(path, ".json") {
		c.loadJSON(path)
	} else if strings.HasSuffix(path, ".yaml") {
		c.loadYAML(path)
	} else {
		panic("unsupported config file format")
	}
}

func (c *Config) loadJSON(path string) {
	// load config from json file
	bts, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	// unmarshal json
	err = json.Unmarshal(bts, c)
	if err != nil {
		panic(err)
	}
}

func (c *Config) loadYAML(path string) {
	bts, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	// unmarshal yaml
	err = yaml.Unmarshal(bts, c)
	if err != nil {
		panic(err)
	}
}
