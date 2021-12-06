package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ConsumerUrl string `mapstructure:"consumerurl"`
	Timeout     int    `mapstructure:"timeout"`
	WorkerCnt   int    `mapstructure:"workercnt"`
	Rates       []int  `mapstructure:"rates"` // list of req/min
	Cnts        []int  `mapstructure:"cnts"`  // list of rate counts
}

func initConfig(filepath string) (Config, error) {
	viper.SetConfigFile(filepath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	config := Config{}
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("read config failed: %s", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshal config failed: %s", err)
	}

	if config.ConsumerUrl == "" {
		return Config{}, fmt.Errorf("config.ConsumerUrl is invalid")
	}
	if config.Timeout <= 0 {
		return Config{}, fmt.Errorf("config.Timeout is invalid")
	}
	if config.WorkerCnt <= 0 {
		return Config{}, fmt.Errorf("config.WorkerCnt is invalid")
	}

	if len(config.Rates) == 0 {
		return Config{}, fmt.Errorf("config.Rates is invalid")
	}
	for _, rate := range config.Rates {
		if rate < 0 {
			return Config{}, fmt.Errorf("config.Rates is invalid")
		}
	}
	if len(config.Rates) != len(config.Cnts) {
		return Config{}, fmt.Errorf("config.Cnts is invalid")
	}
	for _, cnt := range config.Cnts {
		if cnt < 0 {
			return Config{}, fmt.Errorf("config.Cnts is invalid")
		}
	}

	return config, nil
}
