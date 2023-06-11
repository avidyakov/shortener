package config

import (
	"flag"
	"log"
)

type configStruct struct {
	BaseURL    string
	ServerAddr string
}

func newConfig() *configStruct {
	serverAddr := flag.String("a", ":8080", "server address")
	baseURL := flag.String("b", "http://localhost:8080", "base url for short links")
	flag.Parse()

	return &configStruct{
		BaseURL:    *baseURL,
		ServerAddr: *serverAddr,
	}
}

var Config *configStruct

func init() {
	log.Println("Initializing config")
	Config = newConfig()
}
