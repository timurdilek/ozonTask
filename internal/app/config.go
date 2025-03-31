package app

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type PsqlConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
}

type StorageType struct {
	DB string `mapstructure:"DB"`
}
type Config struct {
	Postgres PsqlConfig  `mapstructure:"Postgres"`
	Storage  StorageType `mapstructure:"DB_Type"`
}

const (
	fileName = "config"
	fileType = "yaml"
	filePath = "./config"
)

func LoadConfig() (Config, error) {

	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)
	viper.AddConfigPath(filePath)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, errors.New("error reading config file")
	}

	cfg := Config{}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, errors.New("failed to read config")
	}

	if cfg == (Config{}) {
		return Config{}, errors.New("config is empty")
	}

	return cfg, nil
}

func parseDBConn(cfg PsqlConfig) string {

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)
}
