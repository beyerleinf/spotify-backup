package spotify

import (
	"beyerleinf/spotify-backup/ent"
	"beyerleinf/spotify-backup/internal/server/config"
	http_utils "beyerleinf/spotify-backup/pkg/http"
	"beyerleinf/spotify-backup/pkg/logger"
	util "beyerleinf/spotify-backup/pkg/util"
	"encoding/json"
	"fmt"
)

type SpotifyService struct {
	slogger     *logger.Logger
	db          *ent.Client
	config      *config.Config
	state       string
	redirectUri string
	storageDir  string
}

func New(db *ent.Client, config *config.Config, storageDir string) *SpotifyService {
	return &SpotifyService{
		slogger:     logger.New("spotify", config.Server.LogLevel.Level()),
		db:          db,
		state:       util.GenerateRandomString(16),
		redirectUri: fmt.Sprintf("%s/ui/spotify/callback", config.Spotify.RedirectUri),
		storageDir:  storageDir,
		config:      config,
	}
}

func (s *SpotifyService) GetUserProfile() (SpotifyUserProfile, error) {
	token, err := s.GetAccessToken()
	if err != nil {
		return SpotifyUserProfile{}, err
	}

	headers := map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
	}

	data, _, err := http_utils.Get("https://api.spotify.com/v1/me", headers)
	if err != nil {
		return SpotifyUserProfile{}, err
	}

	var profile SpotifyUserProfile
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return SpotifyUserProfile{}, err
	}

	return profile, nil
}
