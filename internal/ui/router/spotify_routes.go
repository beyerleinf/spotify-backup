package ui

import (
	"beyerleinf/spotify-backup/internal/api/router"
	handler "beyerleinf/spotify-backup/internal/ui/handler"

	"github.com/labstack/echo/v4"
)

func SpotifyRoutes(spotifyHandler *handler.SpotifyHandler) router.RouteGroup {
	return router.RouteGroup{
		Prefix: "/spotify",
		Routes: []router.Route{
			{
				Method:  echo.GET,
				Path:    "/auth",
				Handler: spotifyHandler.SpotifyAuthPage,
			},
			{
				Method:  echo.GET,
				Path:    "/callback",
				Handler: spotifyHandler.SpotifyAuthCallbackPage,
			},
		},
	}
}
