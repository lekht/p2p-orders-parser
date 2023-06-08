package config

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

type Conf struct {
	Asset []string `yaml:"asset"`
	Fiat  []string `yaml:"fiat"`
}

func (p *Conf) Load(path string) error {
	dir, file := filepath.Split(path)
	filename := filepath.Base(path)
	ext := filepath.Ext(file)
	name := filename[0 : len(filename)-len(ext)]

	v := viper.New()
	v.AddConfigPath("." + dir)
	v.SetConfigName(name)
	v.SetConfigType(ext[1:])

	err := v.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "read config")
	}

	err = v.Unmarshal(&p)
	if err != nil {
		return errors.Wrap(err, "config unmarshal")
	}

	return nil
}
