package configuration

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Conf struct {
	SenderType       string
	Interval         int
	HTTPEndpoint     string
	HTTPMethod       string
	EmailAddress     string
	PDAPIKey         string
	PDIntegrationKey string
}

// Reads info from config file
func ReadConfig(fileName string) Conf {
	var configfile = fileName
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Conf
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	//log.Print(config.Index)
	return config
}
