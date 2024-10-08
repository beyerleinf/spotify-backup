package config

import (
	logger "beyerleinf/spotify-backup/pkg/log"
	"fmt"
	goslog "log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig   `mapstructure:"server" env:"SERVER"`
	Database      DatabaseConfig `mapstructure:"database" env:"DB"`
	Spotify       SpotifyConfig  `mapstructure:"spotify" env:"SPOTIFY"`
	EncryptionKey string         `mapstructure:"encryption_key" env:"ENCRYPTION_KEY"`
}

type ServerConfig struct {
	Port     int          `mapstructure:"port" env:"PORT"`
	LogLevel goslog.Level `mapstructure:"loglevel" env:"LOGLEVEL"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host" env:"HOST"`
	Port     int    `mapstructure:"port" env:"PORT"`
	Username string `mapstructure:"username" env:"USERNAME"`
	Password string `mapstructure:"password" env:"PASSWORD"`
	DBName   string `mapstructure:"db_name" env:"NAME"`
}

type SpotifyConfig struct {
	ClientId     string `mapstructure:"client_id" env:"CLIENT_ID"`
	ClientSecret string `mapstructure:"client_secret" env:"CLIENT_SECRET"`
	RedirectUri  string `mapstructure:"redirect_uri" env:"REDIRECT_URI"`
}

var AppConfig Config

func LoadConfig() error {
	slog := logger.New("config", logger.LevelTrace)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.loglevel", 0)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.username", "SpotifyBackup")
	viper.SetDefault("database.password", "secret")
	viper.SetDefault("database.db_name", "SpotifyBackup")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Warn("No config file found. Using environment variables.")
		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}

	slog.Trace("Loaded config", "config", AppConfig)

	return nil
}
