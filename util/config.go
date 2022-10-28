package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var (
	configFile = os.Getenv("CHIAVE_CONFIG_PATH")
)

type Config struct {
	RepFactor        int      `json:"rep_factor"`
	WorkersPerServer int      `json:"workers_per_server"`
	PartitionCount   int      `json:"partition_count"`
	Load             float64  `json:"load"`
	Addresses        []string `json:"addresses"`
}

func LoadConfig() Config {
	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	config := Config{}
	json.Unmarshal([]byte(f), &config)
	return config
}
