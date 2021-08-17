package app

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"wmenjoy/music/etc"
)

func Download() error{
	config, err:= ParseConfig()
	if err != nil {
		return err
	}

	logrus.Printf("%+v", config)
	return nil
}

func ParseConfig() (*etc.Config, error) {
	config := &etc.Config{}

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
