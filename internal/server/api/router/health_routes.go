package router

import (
	"beyerleinf/spotify-backup/internal/server/api/handler"
	"beyerleinf/spotify-backup/pkg/router"

	"github.com/labstack/echo/v4"
)

func HealthRoutes(healthHandler *handler.HealthHandler) router.RouteGroup {
	return router.RouteGroup{
		Prefix: "/health",
		Routes: []router.Route{
			{
				Method:  echo.GET,
				Path:    "",
				Handler: healthHandler.GetHealthStatus,
			},
		},
	}
}
