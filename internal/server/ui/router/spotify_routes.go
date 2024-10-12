package ui

import (
	"beyerleinf/spotify-backup/internal/server/ui/handler"
	"beyerleinf/spotify-backup/pkg/router"

	"github.com/labstack/echo/v4"
)

// SpotifyRoutes returns all routes associated with the /spotify route.
func SpotifyRoutes(spotifyHandler *handler.SpotifyHandler) router.RouteGroup {
	return router.RouteGroup{
		Prefix: "/spotify",
		Routes: []router.Route{
			{
				Method:  echo.GET,
				Path:    "/auth",
				Handler: spotifyHandler.SpotifySettingsPage,
			},
			{
				Method:  echo.GET,
				Path:    "/callback",
				Handler: spotifyHandler.SpotifyAuthCallbackPage,
			},
		},
	}
}
