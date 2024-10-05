package config

import (
	logger "beyerleinf/spotify-backup/pkg/log"
	"fmt"
	goslog "log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port     int
	LogLevel goslog.Level
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

var AppConfig Config

func LoadConfig() error {
	slog := logger.New("config", logger.LevelTrace)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Warn("No config file found. Using environment variables.")
		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	slog.Trace("Read config successfully")

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}

	setDefaults()

	return nil
}

func setDefaults() {
	if AppConfig.Server.Port == 0 {
		AppConfig.Server.Port = 8080
	}

	if AppConfig.Database.Port == 0 {
		AppConfig.Database.Port = 5432 // Default PostgreSQL port
	}

}
