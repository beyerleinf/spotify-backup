package spotify

import (
	"beyerleinf/spotify-backup/ent"
	"beyerleinf/spotify-backup/internal/config"
	logger "beyerleinf/spotify-backup/pkg/log"
	"beyerleinf/spotify-backup/pkg/models"
	util "beyerleinf/spotify-backup/pkg/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		// TODO we are logging this twice, here and where the service is used. Do we need to log here?
		s.slogger.Error("Failed to get access token", "err", err)
		return models.SpotifyUserProfile{}, err
	}

	url := "https://api.spotify.com/v1/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.slogger.Error("Failed to construct user profile request", "err", err)
		return models.SpotifyUserProfile{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.slogger.Error("Failed to get user profile", "err", err)
		return models.SpotifyUserProfile{}, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		s.slogger.Error("Failed to read response data", "err", err)
		return models.SpotifyUserProfile{}, err
	}

	var profile models.SpotifyUserProfile
	err = json.Unmarshal(data, &profile)
	if err != nil {
		s.slogger.Error("Failed to unmarshal response", "err", err)
		return models.SpotifyUserProfile{}, err
	}

	return profile, nil
}
