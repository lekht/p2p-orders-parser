package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Parameters `yaml:"parameters"`
	}

	Parameters struct {
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
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig(): %w", err)
	}

	return cfg, nil
}
