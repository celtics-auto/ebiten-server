package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	AppEnv string `mapstructure:"app_env"`
	Server Server `mapstructure:",squash"`
	Logger Logger `mapstructure:",squash"`
}

type Server struct {
	Port string `mapstructure:"server_port"`
}

type Logger struct {
	Stdout bool `mapstructure:"logger_stdout"`
	File   bool `mapstructure:"logger_file"`
}

func New() (*Config, error) {
	viper.SetDefault("APP_ENV", "production")
	viper.SetDefault("SERVER_PORT", "3000")
	viper.SetDefault("LOGGER_STDOUT", true)
	viper.SetDefault("LOGGER_FILE", false)

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
