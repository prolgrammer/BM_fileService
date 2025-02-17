package config

import (
	"app/config/minio"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type (
	Config struct {
		App   `mapstructure:"app"`
		HTTP  `mapstructure:"http"`
		Minio minio.Minio `mapstructure:"minio"`
	}

	App struct {
		GinMode string `mapstructure:"gin_mode"`
	}

	HTTP struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
)

func New() (*Config, error) {
	cfg := Config{}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config/")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	for _, key := range v.AllKeys() {
		anyValue := v.Get(key)
		str, ok := anyValue.(string)

		if !ok {
			continue
		}

		replaced := os.ExpandEnv(str)
		v.Set(key, replaced)
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling file: %s", err))
	}

	return &cfg, nil
}
