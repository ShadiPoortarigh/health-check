package config

import (
	"encoding/json"
	"os"
)

func ReadConfig(configPath string) (Config, error) {
	var c Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func MustReadConfig(configPath string) Config {
	config, err := ReadConfig(configPath)
	if err != nil {
		panic(err)
	}
	return config
}
