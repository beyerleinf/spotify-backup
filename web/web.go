package web

import (
	"embed"

	"github.com/labstack/echo/v4"
)

//go:embed templates/*
var templates embed.FS
var TemplatesFS = echo.MustSubFS(templates, "templates")

//go:embed static/*
var static embed.FS
var StaticFS = echo.MustSubFS(static, "static")
