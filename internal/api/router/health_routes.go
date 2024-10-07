package router

import (
	"beyerleinf/spotify-backup/internal/api/handler"

	"github.com/labstack/echo/v4"
)

func HealthRoutes(healthHandler *handler.HealthHandler) RouteGroup {
	return RouteGroup{
		Prefix: "/health",
		Routes: []Route{
			{
				Method: echo.GET,
				Path: "",
				Handler: healthHandler.GetHealthStatus,
			},
		},
	}
}