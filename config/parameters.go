package config

import (
	"log"

	"github.com/spf13/viper"
)

type Parameters struct {
	Asset          string   `json:"asset" yaml:"asset"`
	Countries      []string `json:"countries" yaml:"countries"`
	Fiat           string   `json:"fiat" yaml:"fiat"`
	Page           int      `json:"page" yaml:"page"`
	PayTypes       []string `json:"payTypes" yaml:"payTypes"`
	ProMerchantAds bool     `json:"proMerchantAds" yaml:"proMerchantAds"`
	PublisherType  []string `json:"publisherType" yaml:"publisherType"`
	Rows           int      `json:"rows" yaml:"rows"`
	TradeType      string   `json:"tradeType" yaml:"tradeType"`
}

func (p *Parameters) ReqParams() error {
	p.Countries = nil
	p.Page = 1
	p.PayTypes = nil
	p.ProMerchantAds = false
	p.PublisherType = nil
	p.Rows = 10

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
