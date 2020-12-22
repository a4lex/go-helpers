package wrapers

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

//
// Config - struct for config store
//
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DSN      string `yaml:"dsn"`
		MaxIdle  int    `yaml:"maxidle"`
		MaxOpen  int    `yaml:"maxopen"`
	} `yaml:"database"`
}

//
// LoadConfig - Load configuration from file
//
func LoadConfig(config string, cfg *Config) error {
	if c, err := os.Open(config); err == nil {
		defer c.Close()
		decoder := yaml.NewDecoder(c)
		if err = decoder.Decode(&cfg); err != nil {
			return fmt.Errorf("error parsing config file: %v", err)
		}
	} else {
		return fmt.Errorf("error opening config file: %v", err)
	}
	return nil
}
