package config

import "flag"

type ConfigStruct struct {
	BaseURL    string
	ServerAddr string
}

func NewConfig() *ConfigStruct {
	config := new(ConfigStruct)
	flag.StringVar(&config.ServerAddr, "a", ":8080", "server address")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "base url for short links")
	flag.Parse()
	return config
}

var Config *ConfigStruct
