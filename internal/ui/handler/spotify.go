package ui

import (
	"beyerleinf/spotify-backup/internal/config"
	logger "beyerleinf/spotify-backup/pkg/log"
	"beyerleinf/spotify-backup/pkg/service/spotify"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SpotifyHandler struct {
	slogger        *logger.Logger
	spotifyService *spotify.SpotifyService
}

func NewSpotifyHandler(spotifyService *spotify.SpotifyService) *SpotifyHandler {
	return &SpotifyHandler{
		slogger:        logger.New("spotify-ui", config.AppConfig.Server.LogLevel),
		spotifyService: spotifyService,
	}
}

func (s *SpotifyHandler) SpotifyAuthCallbackPage(c echo.Context) error {
	code := c.QueryParams().Get("code")
	state := c.QueryParams().Get("state")

	if code == "" || state == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/ui/spotify/auth?error=code_or_state")
		return nil
	}

	err := s.spotifyService.HandleAuthCallback(code, state)
	if err != nil {
		s.slogger.Error("error handling auth callback", "err", err)
		c.Redirect(http.StatusTemporaryRedirect, "/ui/spotify/auth?error=get_access_token")
		return nil
	}

	c.Redirect(http.StatusTemporaryRedirect, "/ui/spotify/auth")
	return nil
}

func (s *SpotifyHandler) SpotifyAuthPage(c echo.Context) error {
	spotifyAuthUrl := s.spotifyService.GetAuthUrl()
	authError := c.QueryParams().Get("error")

	profile, err := s.spotifyService.GetUserProfile()
	if err != nil {
		s.slogger.Error("Failed to load user profile. Not authenticated?", "err", err)

		return c.Render(http.StatusOK, "spotify_auth", map[string]any{
			"Title":          "Spotify Settings | Spotify Backup",
			"SpotifyAuthUrl": spotifyAuthUrl,
			"HasError":       authError,
		})
	}

	return c.Render(http.StatusOK, "spotify_auth", map[string]any{
		"Title":          "Spotify Settings | Spotify Backup",
		"SpotifyAuthUrl": spotifyAuthUrl,
		"HasError":       authError,
		"Profile":        profile,
	})
}
