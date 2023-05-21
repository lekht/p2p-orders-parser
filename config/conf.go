package config

import (
	"log"

	"github.com/spf13/viper"
)

type Conf struct {
	Asset []string `yaml:"asset"`
	Fiat  []string `yaml:"fiat"`
}

func (p *Conf) ReqParams(path string) error {

	viper.AddConfigPath("./")
	viper.SetConfigName("conf")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("config - ReqParams() - ReadInConfig() error: %s", err)
	}

	err = viper.Unmarshal(&p)
	if err != nil {
		log.Panicf("config - ReqParams() - Unmarshal() error: %s", err)

	}

	return nil
}
