package config

import (
	"github.com/pkg/errors"
	"log"

	"github.com/spf13/viper"
)

type Conf struct {
	Asset []string `yaml:"asset"`
	Fiat  []string `yaml:"fiat"`
}

// todo надо работать с путем, а то он не используется. Моветон прописывать название файла внутри этого метода, потому что его задача открывать и парсить, то что ему скажут. Если проект расширяет
func (p *Conf) Load(path string) error {

	// todo использовать лучше новый объект чтобы полностью управлять его состоянием. А то сейчас мы работаем с какой-то глобальной переменной внутри пакета Viper
	//v:=viper.New()
	//v.AddConfigPath()
	viper.AddConfigPath(".")
	viper.SetConfigName("conf")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	// todo возвращать ошибку, а не панику. Лучше старать ся избегать паник, только если дейтсвительно это необходимо
	// todo прочитать про ошибки в частно wrap, as, is из пакета github.com/pkg/errors
	// todo посмотреть как можно узнать стэк трейс ошибки
	if err != nil {
		return errors.Wrap(err, "read config")
	}

	err = viper.Unmarshal(&p)
	if err != nil {
		log.Panicf("config - Load() - Unmarshal() error: %s", err)

	}

	return nil
}
