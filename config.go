package main

import "os"
import "fmt"
import "encoding/json"
import "net/http"

const (
	defaultconfig = "config.json"
)

type Configuration struct {
	Id      string
	Workers int
	URLs    string
	Port    int
}

func loadConfig(fpath string) Configuration {
	config := Configuration{}
	if &fpath != nil {
		configFile, err := os.Open(fpath)
		if err != nil {
			fmt.Println("Error opening config file:", fpath)
			return config
		}
		decoder := json.NewDecoder(configFile)
		err = decoder.Decode(&config)
		if err != nil {
			fmt.Println("Error parsing config file:", fpath)
		}
	}
	return config
}

func setupServerHandlers() {
	http.HandleFunc("/", rootHandler)
}
