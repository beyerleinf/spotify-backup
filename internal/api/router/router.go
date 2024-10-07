package router

import (
	"github.com/labstack/echo/v4"
)

type Route struct {
	Method      string
	Path        string
	Handler     echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

type RouteGroup struct {
	Prefix      string
	Middlewares []echo.MiddlewareFunc
	Routes      []Route
}

func SetupRoutes(root *echo.Group, groups ...RouteGroup) {
	for _, group := range groups {
		g := root.Group(group.Prefix, group.Middlewares...)

		for _, route := range group.Routes {
			g.Add(route.Method, route.Path, route.Handler, route.Middlewares...)
		}
	}
}
