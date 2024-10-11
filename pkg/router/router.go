package router

import (
	"github.com/labstack/echo/v4"
)

// A RouteGroup is a collection of routes under a common prefix like /foo/bar, /foo/baz, etc.
type RouteGroup struct {
	Prefix      string
	Middlewares []echo.MiddlewareFunc
	Routes      []Route
}

// A Route is a single endpoint like POST /bar and the associated handler and middlewares.
type Route struct {
	Method      string
	Path        string
	Handler     echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

// SetupRoutes adds routes to a echo group
func SetupRoutes(root *echo.Group, groups ...RouteGroup) {
	for _, group := range groups {
		g := root.Group(group.Prefix, group.Middlewares...)

		for _, route := range group.Routes {
			g.Add(route.Method, route.Path, route.Handler, route.Middlewares...)
		}
	}
}
