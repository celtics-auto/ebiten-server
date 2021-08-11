package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Server Server `mapstructure:",squash"`
}

type Server struct {
	Port string `mapstructure:"server_port"`
}

// FIXME: .env not working
func New() (*Config, error) {
	viper.SetDefault("SERVER_PORT", "3000")

	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		zap.L().Error("config: .env file not found")
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
