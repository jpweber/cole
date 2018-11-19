package confiuration

import (
	"log"
	"os"
	"toml"
)

type Config struct {
}

// Reads info from config file
func ReadConfig(fileName string) Config {
	var configfile = fileName
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	//log.Print(config.Index)
	return config
}
