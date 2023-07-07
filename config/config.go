package config

import (
	tracerjaeger "github.com/PatrickChagastavares/game-of-thrones/pkg/tracer/tracer_jaeger"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Env      string               `mapstructure:"env"`
		Port     string               `mapstructure:"port"`
		Tracer   tracerjaeger.Options `mapstructure:"tracer"`
		Database Database             `mapstructure:"database"`
	}
	Database struct {
		Writer string `mapstructure:"writer"`
		Reader string `mapstructure:"reader"`
	}
)

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.Unmarshal(&config)
	return
}
