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
	if code != "" && state != "" {
		s.slogger.Verbose("Spotify Code", "code", code)
		s.spotifyService.GetAuthToken(code, state)
	}

	c.Redirect(http.StatusTemporaryRedirect, "/ui/spotify/auth")
	return nil
}

func (s *SpotifyHandler) SpotifyAuthPage(c echo.Context) error {
	// profile, err := s.spotifyService.GetUserProfile()
	// if err != nil {
	// 	s.slogger.Error("Failed to load user profile. Not authenticated?", "err", err)
	// }

	// s.slogger.Verbose("Profile", "profile", profile)

	spotifyAuthUrl := s.spotifyService.GetAuthUrl()

	return c.Render(http.StatusOK, "spotify_auth", map[string]any{
		"Title":          "Spotify Settings | Spotify Backup",
		"SpotifyAuthUrl": spotifyAuthUrl,
		// "Profile":        profile,
	})
}
