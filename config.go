package main

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type YamlConfig struct {
	Version   string
	Notifiers []map[string]interface{}
}

// LoadConfig reads and parses the config file preliminarily.
// It checks the version of the config and returns map structures for later use.
func LoadConfig(path string) (YamlConfig, error) {
	var config YamlConfig

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(yamlFile), &config)
	if err != nil {
		return config, err
	}

	if config.Version != "1.0" {
		return config, errors.New("Invalid config version")
	}

	return config, nil
}
