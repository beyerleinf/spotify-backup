package main

import (
	"beyerleinf/spotify-backup/internal/api/handler"
	"beyerleinf/spotify-backup/internal/api/router"
	"beyerleinf/spotify-backup/internal/config"
	logger "beyerleinf/spotify-backup/pkg/log"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	slog := logger.New("main", logger.LevelInfo)

	err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load config: %w", err)
		panic(err)
	}

	slog.SetLogLevel(config.AppConfig.Server.LogLevel)

	healthHandler := handler.NewHealthHandler()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(logger.GetEchoLogger())

	api := e.Group("/api")

	router.SetupRoutes(api,
		router.HealthRoutes(healthHandler),
	)

	slog.Info(fmt.Sprintf("Starting server on [::]:%d", config.AppConfig.Server.Port))
	e.Logger.Fatal(e.Start(":8080"))
}
