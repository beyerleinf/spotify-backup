package config

import (
	"beyerleinf/spotify-backup/pkg/logger"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

// Config is the root level configuration struct.
type Config struct {
	Server        ServerConfig   `mapstructure:"server" env:"SERVER"`
	Database      DatabaseConfig `mapstructure:"database" env:"DB"`
	Spotify       SpotifyConfig  `mapstructure:"spotify" env:"SPOTIFY"`
	EncryptionKey string         `mapstructure:"encryption_key" env:"ENCRYPTION_KEY"`
}

// ServerConfig contains setting relating to the http server and the application in general.
type ServerConfig struct {
	Port     int        `mapstructure:"port" env:"PORT"`
	LogLevel slog.Level `mapstructure:"log_level" env:"LOGLEVEL"`
}

// DatabaseConfig contains all database related settings.
type DatabaseConfig struct {
	Host     string `mapstructure:"host" env:"HOST"`
	Port     int    `mapstructure:"port" env:"PORT"`
	Username string `mapstructure:"username" env:"USERNAME"`
	Password string `mapstructure:"password" env:"PASSWORD"`
	DBName   string `mapstructure:"db_name" env:"NAME"`
}

// SpotifyConfig contains all Spotify API related settings.
type SpotifyConfig struct {
	ClientID     string `mapstructure:"client_id" env:"CLIENT_ID"`
	ClientSecret string `mapstructure:"client_secret" env:"CLIENT_SECRET"`
	RedirectURI  string `mapstructure:"redirect_uri" env:"REDIRECT_URI"`
}

// LoadConfig uses viper to load the configuration file.
func LoadConfig() (*Config, error) {
	slogger := logger.New("config", logger.LevelTrace)

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
			slogger.Warn("No config file found. Using environment variables.")
		} else {
			return &Config{}, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return &Config{}, fmt.Errorf("unable to decode into struct: %w", err)
	}

	slogger.Trace("Loaded config", "config", config)

	return &config, nil
}
