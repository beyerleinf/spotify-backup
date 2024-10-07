package spotify

import (
	"beyerleinf/spotify-backup/ent"
	"beyerleinf/spotify-backup/internal/config"
	logger "beyerleinf/spotify-backup/pkg/log"
	"beyerleinf/spotify-backup/pkg/models"
	util "beyerleinf/spotify-backup/pkg/util"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

func (s *SpotifyService) GetAuthUrl() string {
	scope := url.QueryEscape("playlist-read-private user-read-private")

	return fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		config.AppConfig.Spotify.ClientId, scope, url.QueryEscape(s.redirectUri), s.state,
	)
}

// TODO implement token storage in json file with encryption
// TODO also clean this mess up
func (s *SpotifyService) GetAuthToken(authCode string, returnedState string) (string, error) {
	if returnedState != s.state {
		return "", fmt.Errorf("state mismatch")
	}

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", authCode)
	form.Add("redirect_uri", s.redirectUri)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
	if err != nil {
		s.slogger.Error("Failed to construct user profile request", "err", err)
		return "", err
	}

	clientIdAndSecret := fmt.Sprintf("%s:%s", config.AppConfig.Spotify.ClientId, config.AppConfig.Spotify.ClientSecret)
	authHeaderValue := base64.StdEncoding.EncodeToString([]byte(clientIdAndSecret))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", authHeaderValue))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.slogger.Error("Failed to get user profile", "err", err)
		return "", err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		s.slogger.Error("Failed to read response data", "err", err)
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		s.slogger.Error("Token request failed", "status", res.Status, "body", string(data))
		return "", fmt.Errorf("token request failed: %s - %s", res.Status, string(data))
	}

	var tokenResponse models.AuthTokenResponse
	err = json.Unmarshal(data, &tokenResponse)
	if err != nil {
		s.slogger.Error("Failed to unmarshal response", "err", err)
		return "", err
	}

	s.slogger.Verbose("token res", "token", tokenResponse)

	return tokenResponse.AccessToken, nil
}

func (s *SpotifyService) GetUserProfile() (models.SpotifyUserProfile, error) {

	url := "https://api.spotify.com/v1/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.slogger.Error("Failed to construct user profile request", "err", err)
		return models.SpotifyUserProfile{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "token"))
	s.slogger.Info("after req construct")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.slogger.Error("Failed to get user profile", "err", err)
		return models.SpotifyUserProfile{}, err
	}

	s.slogger.Info("after req")

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		s.slogger.Error("Failed to read response data", "err", err)
		return models.SpotifyUserProfile{}, err
	}

	s.slogger.Info("after read", "data", data)

	var profile models.SpotifyUserProfile
	err = json.Unmarshal(data, &profile)
	if err != nil {
		s.slogger.Error("Failed to unmarshal response", "err", err)
		return models.SpotifyUserProfile{}, err
	}

	s.slogger.Info("Profile Name", "name", profile.DisplayName)

	return profile, nil
}
