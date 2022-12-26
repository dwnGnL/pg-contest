package config

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func FromFile(filename string) *Config {
	conf, err := ViperInitConfig(filename)
	if err != nil {
		logrus.WithError(err).Fatal("cannot load config")
	}

	return conf
}

func ViperInitConfig(configFile string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetConfigFile(configFile)
	v.SetEnvPrefix("contests")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = v.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
