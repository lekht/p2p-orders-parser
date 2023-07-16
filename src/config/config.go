package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

type Conf struct {
	Asset []string `yaml:"asset"`
	Fiat  []string `yaml:"fiat"`
}

func (p *Conf) Parse(path string) error {
	if path == "" {
		files, err := os.ReadDir(".")
		if err != nil {
			return errors.Wrap(err, "failed to find config")
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			filename := file.Name()

			if ext := filepath.Ext(filename); ext != ".yml" {
				continue
			}

			if err := p.Load("/" + filename); err != nil {
				return errors.Wrap(err, "failed to load config")
			}
			return nil
		}

	}

	if err := p.Load(path); err != nil {
		return errors.Wrap(err, "failed to load config")
	}
	return nil

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
		return errors.Wrap(err, "failed to read yaml config")
	}

	err = v.Unmarshal(&p)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config to struct")
	}

	return nil
}
