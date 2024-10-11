package handler

import (
	"beyerleinf/spotify-backup/internal/server/config"
	"beyerleinf/spotify-backup/pkg/logger"
	"beyerleinf/spotify-backup/pkg/service/spotify"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SpotifyHandler struct {
	slogger        *logger.Logger
	spotifyService *spotify.SpotifyService
	config         *config.Config
}

const pageTitle = "Spotify Settings | Spotify Backup"

func NewSpotifyHandler(spotifyService *spotify.SpotifyService, config *config.Config) *SpotifyHandler {
	return &SpotifyHandler{
		slogger:        logger.New("spotify-ui", config.Server.LogLevel),
		spotifyService: spotifyService,
		config:         config,
	}
}

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

func (s *SpotifyHandler) SpotifyAuthPage(c echo.Context) error {
	spotifyAuthUrl := s.spotifyService.GetAuthUrl()
	authError := c.QueryParams().Get("error")

	profile, err := s.spotifyService.GetUserProfile()
	if err != nil {
		s.slogger.Error("Failed to load user profile. Not authenticated?", "err", err)

		return c.Render(http.StatusOK, "spotify_auth", map[string]any{
			"Title":          pageTitle,
			"SpotifyAuthUrl": spotifyAuthUrl,
			"HasError":       authError,
		})
	}

	return c.Render(http.StatusOK, "spotify_auth", map[string]any{
		"Title":          pageTitle,
		"SpotifyAuthUrl": spotifyAuthUrl,
		"HasError":       authError,
		"Profile":        profile,
	})
}
