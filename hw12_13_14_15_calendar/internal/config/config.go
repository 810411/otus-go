package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	HTTP struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"http"`
	Log struct {
		File  string `json:"file"`
		Level string `json:"level"`
	} `json:"log"`
	Repository struct {
		Type string `json:"type"`
		Dsn  string `json:"dsn"`
	} `json:"repository"`
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c := &Config{}
	dec := json.NewDecoder(file)
	if err = dec.Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}
