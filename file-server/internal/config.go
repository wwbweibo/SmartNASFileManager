package internal

import (
	"encoding/json"
	"fileserver/internal/adapters/dl"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var config = &Config{}

type Config struct {
	// StartWithBackendScan is a flag to start with backend scan, if true, the program will start with backend scan
	StartWithBackendScan bool   `json:"start_with_backend_scan" yaml:"start_with_backend_scan"`
	NasRootPath          string `json:"nas_root_path" yaml:"nas_root_path"`
	DBPath               string `json:"db_path" yaml:"db_path"`
	CachePath            string `json:"cache_path" yaml:"cache_path"`
	ScanOption           struct {
		Path       []string `json:"path" yaml:"path"`
		RegexPath  []string `json:"regex_path" yaml:"regex_path"`
		Extensions []string `json:"extensions" yaml:"extensions"`
	} `json:"scan_option" yaml:"scan_option"`
	DLConfiguration dl.Config `json:"dl_configuration" yaml:"dl_configuration"`
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

func GetConfig() *Config {
	return config
}
