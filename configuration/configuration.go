package configuration

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Conf struct {
	SenderType       string `env:"SENDER_TYPE,required"`
	Interval         int    `env:"INTERVAL" envDefault:"65"`
	HTTPEndpoint     string `env:"HTTP_ENDPOINT"`
	HTTPMethod       string `env:"HTTP_METHOD" envDefault:"POST"`
	EmailAddress     string `env:"EMAIL_ADDR"`
	PDIntegrationKey string `env:"PD_KEY"`
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
