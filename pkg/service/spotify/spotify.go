package spotify

import (
	"beyerleinf/spotify-backup/internal/server/config"
	"beyerleinf/spotify-backup/pkg/logger"
	"beyerleinf/spotify-backup/pkg/request"
	util "beyerleinf/spotify-backup/pkg/util"
	"context"
	"encoding/json"
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

// GetUserProfile returns a [UserProfile] from Spotify's API.
// [Get User Profile API]: https://developer.spotify.com/documentation/web-api/reference/get-current-users-profile
func (s *Service) GetUserProfile() (UserProfile, error) {
	ctx := context.Background()

	token, err := s.GetAccessToken()
	if err != nil {
		return UserProfile{}, err
	}

	headers := map[string][]string{
		"Authorization": {"Bearer " + token},
	}

	data, _, err := request.Get(ctx, "https://api.spotify.com/v1/me", headers)
	if err != nil {
		return UserProfile{}, err
	}

	var profile UserProfile
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}
