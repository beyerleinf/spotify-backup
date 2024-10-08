package spotify

import (
	"beyerleinf/spotify-backup/ent"
	"beyerleinf/spotify-backup/internal/config"
	http_utils "beyerleinf/spotify-backup/pkg/http"
	logger "beyerleinf/spotify-backup/pkg/log"
	"beyerleinf/spotify-backup/pkg/models"
	util "beyerleinf/spotify-backup/pkg/util"
	"encoding/json"
	"fmt"
)

type SpotifyService struct {
	slogger     *logger.Logger
	db          *ent.Client
	state       string
	redirectUri string
}

func New(db *ent.Client) *SpotifyService {
	return &SpotifyService{
		slogger:     logger.New("spotify", config.AppConfig.Server.LogLevel.Level()),
		db:          db,
		state:       util.GenerateRandomString(16),
		redirectUri: fmt.Sprintf("%s/ui/spotify/callback", config.AppConfig.Spotify.RedirectUri),
	}
}

func (s *SpotifyService) GetUserProfile() (models.SpotifyUserProfile, error) {
	token, err := s.GetAccessToken()
	if err != nil {
		return models.SpotifyUserProfile{}, err
	}

	headers := map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
	}

	data, _, err := http_utils.Get("https://api.spotify.com/v1/me", headers)
	if err != nil {
		return models.SpotifyUserProfile{}, err
	}

	var profile models.SpotifyUserProfile
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return models.SpotifyUserProfile{}, err
	}

	return profile, nil
}
