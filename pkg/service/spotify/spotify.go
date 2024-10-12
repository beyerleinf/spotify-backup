package spotify

import (
	"beyerleinf/spotify-backup/internal/server/config"
	"beyerleinf/spotify-backup/pkg/logger"
	util "beyerleinf/spotify-backup/pkg/util"
)

// A Service instance.
type Service struct {
	slogger     *logger.Logger
	config      *config.Config
	state       string
	redirectURI string
	storageDir  string
}

// New creates a [Service] instance.
func New(config *config.Config, storageDir string) *Service {
	return &Service{
		slogger:     logger.New("spotify", config.Server.LogLevel.Level()),
		state:       util.GenerateRandomString(16),
		redirectURI: config.Spotify.RedirectURI + "/ui/spotify/callback",
		storageDir:  storageDir,
		config:      config,
	}
}
