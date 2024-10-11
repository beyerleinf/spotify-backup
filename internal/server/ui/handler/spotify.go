package handler

import (
	"beyerleinf/spotify-backup/internal/server/config"
	"beyerleinf/spotify-backup/pkg/logger"
	"beyerleinf/spotify-backup/pkg/service/spotify"
	"net/http"

	"github.com/labstack/echo/v4"
)

// A SpotifyHandler instance.
type SpotifyHandler struct {
	slogger        *logger.Logger
	spotifyService *spotify.Service
	config         *config.Config
}

const pageTitle = "Spotify Settings | Spotify Backup"

// NewSpotifyHandler creates a new instance.
func NewSpotifyHandler(spotifyService *spotify.Service, config *config.Config) *SpotifyHandler {
	return &SpotifyHandler{
		slogger:        logger.New("spotify-ui", config.Server.LogLevel),
		spotifyService: spotifyService,
		config:         config,
	}
}

// SpotifyAuthCallbackPage handles callback requests from Spotify's Authentication API.
func (s *SpotifyHandler) SpotifyAuthCallbackPage(c echo.Context) error {
	code := c.QueryParams().Get("code")
	state := c.QueryParams().Get("state")

	if code == "" || state == "" {
		err := c.Redirect(http.StatusTemporaryRedirect, "/ui/spotify/auth?error=code_or_state")
		if err != nil {
			return err
		}

		return nil
	}

	err := s.spotifyService.HandleAuthCallback(code, state)
	if err != nil {
		s.slogger.Error("error handling auth callback", "err", err)

		err = c.Redirect(http.StatusTemporaryRedirect, "/ui/spotify/auth?error=get_access_token")
		if err != nil {
			return err
		}

		return nil
	}

	err = c.Redirect(http.StatusTemporaryRedirect, "/ui/spotify/auth")
	if err != nil {
		return err
	}

	return nil
}

// SpotifySettingsPage serves the Spotify settings page.
func (s *SpotifyHandler) SpotifySettingsPage(c echo.Context) error {
	const templateName = "spotify_settings"

	authURL := s.spotifyService.GetAuthURL()
	authError := c.QueryParams().Get("error")

	profile, err := s.spotifyService.GetUserProfile()
	if err != nil {
		s.slogger.Error("Failed to load user profile. Not authenticated?", "err", err)

		return c.Render(http.StatusOK, templateName, map[string]any{
			"Title":    pageTitle,
			"AuthURL":  authURL,
			"HasError": authError,
		})
	}

	return c.Render(http.StatusOK, templateName, map[string]any{
		"Title":    pageTitle,
		"AuthURL":  authURL,
		"HasError": authError,
		"Profile":  profile,
	})
}
