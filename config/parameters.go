package config

import (
	"log"

	"github.com/spf13/viper"
)

type Parameters struct {
	Asset []string `yaml:"asset"`
	Fiat  []string `yaml:"fiat"`
}

func (p *Parameters) ReqParams() error {

	viper.AddConfigPath("./config/")
	viper.SetConfigName("parameters")
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
